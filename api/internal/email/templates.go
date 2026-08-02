package email

import (
	"fmt"
	"strings"
)

type emailView struct {
	LogoURL         string
	CustomerName    string
	OrderNo         string
	Status          string
	PaymentProvider string
	PaymentRef      string
	PaidAt          string
	Email           string
	Phone           string
	Address         string
	DeliveryZone    string
	Subtotal        string
	DeliveryFee     string
	Total           string
	Items           []emailItemView
}

type emailItemView struct {
	Name     string
	Size     string
	Quantity int
	Price    string
	Line     string
}

// newEmailView flattens an order into template-safe display values.
func newEmailView(siteURL string, order Order) emailView {
	items := make([]emailItemView, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, emailItemView{
			Name:     item.Name,
			Size:     item.Size,
			Quantity: item.Quantity,
			Price:    formatNaira(item.PriceKobo),
			Line:     formatNaira(item.LineKobo),
		})
	}

	return emailView{
		LogoURL:         strings.TrimRight(siteURL, "/") + "/images/brand/oyecouture-logo.png",
		CustomerName:    order.CustomerName,
		OrderNo:         order.OrderNo,
		Status:          "Paid",
		PaymentProvider: displayProvider(order.Payment.Provider),
		PaymentRef:      order.Payment.Reference,
		PaidAt:          order.Payment.PaidAt,
		Email:           order.Email,
		Phone:           order.Phone,
		Address:         order.Address,
		DeliveryZone:    order.DeliveryZone,
		Subtotal:        formatNaira(order.SubtotalKobo),
		DeliveryFee:     formatNaira(order.DeliveryFeeKobo),
		Total:           formatNaira(order.TotalKobo),
		Items:           items,
	}
}

// displayProvider formats known provider identifiers for customer emails.
func displayProvider(provider string) string {
	switch strings.ToLower(provider) {
	case "paystack":
		return "Paystack"
	case "sendbyte":
		return "SendByte"
	default:
		return provider
	}
}

// renderText creates a plain-text fallback for clients that do not render HTML.
func renderText(emailType string, view emailView) string {
	var builder strings.Builder
	if emailType == TypePaymentReceipt {
		builder.WriteString("Payment received for " + view.OrderNo + "\n\n")
		builder.WriteString("Hi " + view.CustomerName + ",\n")
		builder.WriteString("We received your payment of " + view.Total + ".\n\n")
		builder.WriteString("Payment reference: " + view.PaymentRef + "\n")
		builder.WriteString("Payment provider: " + view.PaymentProvider + "\n")
		builder.WriteString("Paid at: " + view.PaidAt + "\n\n")
	} else {
		builder.WriteString("Order confirmed - " + view.OrderNo + "\n\n")
		builder.WriteString("Hi " + view.CustomerName + ",\n")
		builder.WriteString("Your Oyecouture order has been confirmed and payment received.\n\n")
	}

	builder.WriteString("Items:\n")
	for _, item := range view.Items {
		builder.WriteString("- " + item.Name)
		if item.Size != "" {
			builder.WriteString(" (" + item.Size + ")")
		}
		builder.WriteString(fmt.Sprintf(" x%d - %s\n", item.Quantity, item.Line))
	}
	builder.WriteString("\nSubtotal: " + view.Subtotal + "\n")
	builder.WriteString("Delivery: " + view.DeliveryFee + "\n")
	builder.WriteString("Total: " + view.Total + "\n\n")
	builder.WriteString("Delivery details:\n")
	builder.WriteString(view.CustomerName + "\n")
	builder.WriteString(view.Phone + "\n")
	builder.WriteString(view.Address + "\n")
	builder.WriteString(view.DeliveryZone + "\n")

	return builder.String()
}

// formatNaira converts integer kobo into a naira display amount.
func formatNaira(amountKobo int64) string {
	naira := amountKobo / 100
	return "₦" + commaInt(naira)
}

// commaInt formats a positive integer with thousands separators.
func commaInt(value int64) string {
	valueString := fmt.Sprintf("%d", value)
	if len(valueString) <= 3 {
		return valueString
	}

	var parts []string
	for len(valueString) > 3 {
		parts = append([]string{valueString[len(valueString)-3:]}, parts...)
		valueString = valueString[:len(valueString)-3]
	}
	parts = append([]string{valueString}, parts...)

	return strings.Join(parts, ",")
}

const emailBaseStyles = `
  body { margin: 0; padding: 0; background: #f4f0ec; color: #24211f; font-family: Arial, Helvetica, sans-serif; }
  .preheader { display: none; max-height: 0; overflow: hidden; opacity: 0; color: transparent; }
  .shell { width: 100%; background: #f4f0ec; padding: 36px 14px; }
  .card { max-width: 660px; margin: 0 auto; background: #fffdf9; border: 1px solid #e5dccf; }
  .topline { height: 6px; background: #24211f; }
  .header { padding: 34px 38px 30px; border-bottom: 1px solid #e8dfd3; }
  .brand { color: #24211f; font-size: 18px; font-weight: 700; letter-spacing: 5px; text-transform: uppercase; }
  .brand span { display: block; margin-top: 5px; color: #b49a5e; font-size: 10px; font-weight: 600; letter-spacing: 6px; }
  .eyebrow { margin: 34px 0 10px; color: #b49a5e; font-size: 12px; font-weight: 700; letter-spacing: 4px; text-transform: uppercase; }
  h1 { margin: 0; color: #24211f; font-size: 31px; line-height: 1.18; font-weight: 700; }
  p { color: #5d5751; font-size: 15px; line-height: 1.75; margin: 0; }
  .summary { width: 100%; margin-top: 24px; border-collapse: collapse; }
  .summary td { padding: 14px 16px; background: #f7f3ee; border: 1px solid #e8dfd3; color: #4a443e; font-size: 13px; vertical-align: top; }
  .summary strong { display: block; margin-top: 5px; color: #24211f; font-size: 17px; }
  .badge { display: inline-block; padding: 6px 12px; border: 1px solid #79a486; color: #457457; border-radius: 999px; font-size: 12px; font-weight: 700; }
  .section { padding: 28px 38px; border-bottom: 1px solid #e8dfd3; }
  .section.alt { background: #fbf8f3; }
  .label { color: #9a9188; font-size: 11px; font-weight: 700; letter-spacing: 3px; text-transform: uppercase; margin-bottom: 16px; }
  .row { width: 100%; border-collapse: collapse; }
  .row td { padding: 8px 0; color: #4a443e; font-size: 14px; line-height: 1.45; vertical-align: top; }
  .row .value { color: #24211f; text-align: right; font-weight: 700; word-break: break-word; }
  .items { width: 100%; border-collapse: collapse; }
  .items th { padding: 0 0 11px; color: #9a9188; font-size: 11px; font-weight: 700; letter-spacing: 2.4px; text-transform: uppercase; text-align: left; }
  .items td { padding: 16px 0; border-top: 1px solid #e8dfd3; color: #24211f; font-size: 14px; line-height: 1.45; vertical-align: top; }
  .items .muted { display: block; margin-top: 5px; color: #8d847c; font-size: 13px; }
  .money { text-align: right; font-weight: 700; white-space: nowrap; }
  .total { font-size: 19px; }
  .grand-total td { padding-top: 15px; border-top: 1px solid #d7cabc; }
  .note { padding: 18px 20px; background: #f7f3ee; border-left: 3px solid #b49a5e; color: #5d5751; font-size: 14px; line-height: 1.65; }
  .footer { padding: 26px 38px 34px; color: #8d847c; font-size: 13px; line-height: 1.7; }
  .footer strong { color: #24211f; letter-spacing: 2px; text-transform: uppercase; }
  @media only screen and (max-width: 600px) {
    .shell { padding: 18px 8px; }
    .header, .section, .footer { padding-left: 22px; padding-right: 22px; }
    h1 { font-size: 25px; }
    .summary td { display: block; width: auto; }
  }
`

const orderConfirmationHTML = `<!doctype html>
<html>
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><style>` + emailBaseStyles + `</style></head>
<body>
  <div class="preheader">Your Oyecouture order {{.OrderNo}} is confirmed and payment has been received.</div>
  <div class="shell">
    <div class="card">
      <div class="topline"></div>
      <div class="header">
        <div class="brand">Oyecouture<span>Bespoke</span></div>
        <div class="eyebrow">Order confirmed</div>
        <h1>Thank you, {{.CustomerName}}.</h1>
        <table class="summary">
          <tr>
            <td>Order<strong>{{.OrderNo}}</strong></td>
            <td>Status<br><span class="badge">{{.Status}}</span></td>
            <td>Total<strong>{{.Total}}</strong></td>
          </tr>
        </table>
      </div>
      <div class="section">
        <p>Your Oyecouture order has been confirmed and your payment has been received. We will prepare your pieces with care and keep you updated as fulfilment begins.</p>
      </div>
      <div class="section">
        <div class="label">Items</div>
        <table class="items">
          <tr><th>Piece</th><th style="text-align:right;">Amount</th></tr>
          {{range .Items}}
          <tr>
            <td>{{.Name}}<span class="muted">{{if .Size}}Size {{.Size}} · {{end}}Qty {{.Quantity}}</span></td>
            <td class="money">{{.Line}}</td>
          </tr>
          {{end}}
        </table>
      </div>
      <div class="section alt">
        <div class="label">Delivery details</div>
        <p>{{.CustomerName}}<br>{{.Phone}}<br>{{.Email}}<br>{{.Address}}<br>{{.DeliveryZone}}</p>
      </div>
      <div class="section">
        <table class="row">
          <tr><td>Subtotal</td><td class="value">{{.Subtotal}}</td></tr>
          <tr><td>Delivery</td><td class="value">{{.DeliveryFee}}</td></tr>
          <tr class="grand-total"><td>Total</td><td class="value total">{{.Total}}</td></tr>
        </table>
      </div>
      <div class="section">
        <div class="note">Your order is now in our fulfilment queue. If anything needs your attention, the Oyecouture team will contact you using the details above.</div>
      </div>
      <div class="footer"><strong>Oyecouture Bespoke</strong><br>This confirmation was sent after secure payment verification.</div>
    </div>
  </div>
</body>
</html>`

const paymentReceiptHTML = `<!doctype html>
<html>
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><style>` + emailBaseStyles + `</style></head>
<body>
  <div class="preheader">Payment of {{.Total}} received for Oyecouture order {{.OrderNo}}.</div>
  <div class="shell">
    <div class="card">
      <div class="topline"></div>
      <div class="header">
        <div class="brand">Oyecouture<span>Bespoke</span></div>
        <div class="eyebrow">Payment receipt</div>
        <h1>Payment received.</h1>
        <table class="summary">
          <tr>
            <td>Amount paid<strong>{{.Total}}</strong></td>
            <td>Order<strong>{{.OrderNo}}</strong></td>
            <td>Status<br><span class="badge">{{.Status}}</span></td>
          </tr>
        </table>
      </div>
      <div class="section">
        <p>Hi {{.CustomerName}}, we received your payment for order {{.OrderNo}}. Keep this receipt for your records.</p>
      </div>
      <div class="section alt">
        <div class="label">Payment details</div>
        <table class="row">
          <tr><td>Amount paid</td><td class="value total">{{.Total}}</td></tr>
          <tr><td>Provider</td><td class="value">{{.PaymentProvider}}</td></tr>
          <tr><td>Reference</td><td class="value">{{.PaymentRef}}</td></tr>
          <tr><td>Paid at</td><td class="value">{{.PaidAt}}</td></tr>
        </table>
      </div>
      <div class="section">
        <div class="label">Order summary</div>
        <table class="items">
          <tr><th>Piece</th><th style="text-align:right;">Amount</th></tr>
          {{range .Items}}
          <tr>
            <td>{{.Name}}<span class="muted">{{if .Size}}Size {{.Size}} · {{end}}Qty {{.Quantity}}</span></td>
            <td class="money">{{.Line}}</td>
          </tr>
          {{end}}
        </table>
      </div>
      <div class="section">
        <table class="row">
          <tr><td>Subtotal</td><td class="value">{{.Subtotal}}</td></tr>
          <tr><td>Delivery</td><td class="value">{{.DeliveryFee}}</td></tr>
          <tr class="grand-total"><td>Total paid</td><td class="value total">{{.Total}}</td></tr>
        </table>
      </div>
      <div class="footer"><strong>Oyecouture Bespoke</strong><br>If you need help, reply to this email with your order number and payment reference.</div>
    </div>
  </div>
</body>
</html>`
