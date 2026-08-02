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
  body { margin: 0; padding: 0; background: #f7f4ef; color: #24211f; font-family: Arial, Helvetica, sans-serif; }
  .shell { width: 100%; background: #f7f4ef; padding: 32px 12px; }
  .card { max-width: 640px; margin: 0 auto; background: #ffffff; border: 1px solid #e8e0d6; }
  .header { padding: 30px 34px 24px; border-bottom: 1px solid #ece5dc; }
  .logo { width: 168px; max-width: 70%; height: auto; display: block; }
  .eyebrow { margin: 28px 0 8px; color: #b49a5e; font-size: 12px; letter-spacing: 3px; text-transform: uppercase; }
  h1 { margin: 0; color: #24211f; font-size: 28px; line-height: 1.2; font-weight: 600; }
  p { color: #5d5751; font-size: 15px; line-height: 1.7; margin: 0; }
  .section { padding: 26px 34px; border-bottom: 1px solid #ece5dc; }
  .label { color: #9a9188; font-size: 11px; letter-spacing: 2.4px; text-transform: uppercase; margin-bottom: 12px; }
  .row { width: 100%; border-collapse: collapse; }
  .row td { padding: 9px 0; color: #4a443e; font-size: 14px; vertical-align: top; }
  .row .value { color: #24211f; text-align: right; font-weight: 600; }
  .items { width: 100%; border-collapse: collapse; }
  .items th { padding: 0 0 10px; color: #9a9188; font-size: 11px; letter-spacing: 2px; text-transform: uppercase; text-align: left; }
  .items td { padding: 14px 0; border-top: 1px solid #ece5dc; color: #24211f; font-size: 14px; vertical-align: top; }
  .items .muted { display: block; margin-top: 4px; color: #948c84; font-size: 13px; }
  .total { font-size: 18px; }
  .footer { padding: 24px 34px 32px; color: #8d847c; font-size: 13px; line-height: 1.7; }
`

const orderConfirmationHTML = `<!doctype html>
<html>
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><style>` + emailBaseStyles + `</style></head>
<body>
  <div class="shell">
    <div class="card">
      <div class="header">
        <img class="logo" src="{{.LogoURL}}" alt="Oyecouture">
        <div class="eyebrow">Order confirmed</div>
        <h1>Thank you, {{.CustomerName}}.</h1>
      </div>
      <div class="section">
        <p>Your Oyecouture order has been confirmed and your payment has been received. We will prepare your pieces and keep you updated as fulfillment begins.</p>
      </div>
      <div class="section">
        <div class="label">Order details</div>
        <table class="row">
          <tr><td>Order number</td><td class="value">{{.OrderNo}}</td></tr>
          <tr><td>Status</td><td class="value">{{.Status}}</td></tr>
          <tr><td>Total</td><td class="value total">{{.Total}}</td></tr>
        </table>
      </div>
      <div class="section">
        <div class="label">Items</div>
        <table class="items">
          <tr><th>Piece</th><th style="text-align:right;">Amount</th></tr>
          {{range .Items}}
          <tr>
            <td>{{.Name}}<span class="muted">{{if .Size}}Size {{.Size}} · {{end}}Qty {{.Quantity}}</span></td>
            <td style="text-align:right;">{{.Line}}</td>
          </tr>
          {{end}}
        </table>
      </div>
      <div class="section">
        <div class="label">Delivery details</div>
        <p>{{.CustomerName}}<br>{{.Phone}}<br>{{.Email}}<br>{{.Address}}<br>{{.DeliveryZone}}</p>
      </div>
      <div class="section">
        <table class="row">
          <tr><td>Subtotal</td><td class="value">{{.Subtotal}}</td></tr>
          <tr><td>Delivery</td><td class="value">{{.DeliveryFee}}</td></tr>
          <tr><td>Total</td><td class="value total">{{.Total}}</td></tr>
        </table>
      </div>
      <div class="footer">Oyecouture Bespoke<br>This confirmation was sent after server-side payment verification.</div>
    </div>
  </div>
</body>
</html>`

const paymentReceiptHTML = `<!doctype html>
<html>
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><style>` + emailBaseStyles + `</style></head>
<body>
  <div class="shell">
    <div class="card">
      <div class="header">
        <img class="logo" src="{{.LogoURL}}" alt="Oyecouture">
        <div class="eyebrow">Payment receipt</div>
        <h1>Payment received.</h1>
      </div>
      <div class="section">
        <p>Hi {{.CustomerName}}, we received your payment of {{.Total}} for order {{.OrderNo}}.</p>
      </div>
      <div class="section">
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
            <td style="text-align:right;">{{.Line}}</td>
          </tr>
          {{end}}
        </table>
      </div>
      <div class="section">
        <table class="row">
          <tr><td>Subtotal</td><td class="value">{{.Subtotal}}</td></tr>
          <tr><td>Delivery</td><td class="value">{{.DeliveryFee}}</td></tr>
          <tr><td>Total</td><td class="value total">{{.Total}}</td></tr>
        </table>
      </div>
      <div class="footer">Keep this receipt for your records. If you need help, reply to this email with your order number.</div>
    </div>
  </div>
</body>
</html>`
