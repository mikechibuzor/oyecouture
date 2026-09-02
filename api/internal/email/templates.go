package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"strings"
)

// firstName keeps the greeting personal: "Thank you, Adeola" reads as a person
// writing where the full legal name reads as a system generating. Falls back to
// the whole string so a single-word name is still greeted rather than dropped.
func firstName(fullName string) string {
	trimmed := strings.TrimSpace(fullName)
	if first, _, found := strings.Cut(trimmed, " "); found && first != "" {
		return first
	}

	return trimmed
}

type emailView struct {
	LogoURL         string
	CustomerName    string
	FirstName       string
	OrderURL        string
	OrderDate       string
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

type fulfillmentEmailView struct {
	CustomerName          string
	FirstName             string
	OrderNo               string
	StageLabel            string
	Eyebrow               string
	Heading               string
	Intro                 string
	PublicNote            string
	TrackingURL           string
	CarrierName           string
	CarrierTrackingNumber string
	CarrierTrackingURL    string
	EstimatedDeliveryDate string
	Items                 []emailItemView
	Subtotal              string
	DeliveryFee           string
	Total                 string
	Progress              []progressStepView
	SignedForBy           string
	SignedAt              string
	DeliveryLocation      string
	RiderName             string
	RiderContact          string
	DeliveryWindow        string
}

// progressStepView drives the four-step bar in the client's design. The API
// tracks eight internal stages; a customer is shown these four, so workshop
// states like quality check stay private rather than reading as delay.
type progressStepView struct {
	Label string
	Done  bool
}

// progressFor maps an internal stage onto the customer-facing four.
func progressFor(emailType string) []progressStepView {
	reached := 2
	switch emailType {
	case TypeOrderDispatched, TypeOutForDelivery:
		// Out for delivery is still the shipped step: the parcel has not
		// arrived, so lighting the delivered marker would contradict the
		// message beside it.
		reached = 3
	case TypeOrderDelivered:
		reached = 4
	}

	labels := []string{"Confirmed", "In production", "Shipped", "Delivered"}
	steps := make([]progressStepView, 0, len(labels))
	for index, label := range labels {
		steps = append(steps, progressStepView{Label: label, Done: index < reached})
	}

	return steps
}

// renderFulfillmentEmail creates branded HTML and plain-text milestone updates.
func renderFulfillmentEmail(
	emailType string,
	payload FulfillmentPayload,
) (string, string, error) {
	items := make([]emailItemView, 0, len(payload.Items))
	for _, item := range payload.Items {
		items = append(items, emailItemView{
			Name:     item.Name,
			Size:     item.Size,
			Quantity: item.Quantity,
			Line:     formatNaira(item.LineKobo),
		})
	}

	view := fulfillmentEmailView{
		CustomerName:          payload.CustomerName,
		FirstName:             firstName(payload.CustomerName),
		OrderNo:               payload.OrderNo,
		StageLabel:            payload.StageLabel,
		PublicNote:            payload.PublicNote,
		TrackingURL:           payload.TrackingURL,
		CarrierName:           payload.CarrierName,
		CarrierTrackingNumber: payload.CarrierTrackingNumber,
		CarrierTrackingURL:    payload.CarrierTrackingURL,
		EstimatedDeliveryDate: payload.EstimatedDeliveryDate,
		Items:                 items,
		Subtotal:              formatNaira(payload.SubtotalKobo),
		DeliveryFee:           formatNaira(payload.DeliveryFeeKobo),
		Total:                 formatNaira(payload.TotalKobo),
		Progress:              progressFor(emailType),
		SignedForBy:           payload.SignedForBy,
		SignedAt:              payload.SignedAt,
		DeliveryLocation:      payload.DeliveryLocation,
		RiderName:             payload.RiderName,
		RiderContact:          payload.RiderContact,
		DeliveryWindow:        payload.DeliveryWindow,
	}

	htmlTemplate := fulfillmentUpdateHTML

	switch emailType {
	case TypeOrderDispatched:
		// Design 02, Shipped, has its own layout: progress bar, a tracking
		// panel and the order reproduced. It gets its own template rather than
		// bending the generic one with conditionals.
		// The heading lives in that template, not here: it interpolates the
		// customer's name, and building markup around untrusted input in Go
		// would mean marking it template.HTML and losing the escaping.
		htmlTemplate = orderDispatchedHTML
		view.Eyebrow = "On its way"
		view.Intro = "Everything has been pressed, checked and packed by hand. Your parcel is now with the courier."
	case TypeOutForDelivery:
		// Design 03, Out for delivery: the rider panel in place of tracking.
		htmlTemplate = outForDeliveryHTML
		view.Eyebrow = "Arriving today"
		view.Intro = "Your order is on the road and should reach you today. Please keep your phone nearby, the rider will call on arrival. If today is inconvenient, reply and we will arrange another slot at no cost."
	case TypeOrderDelivered:
		// Design 04, Delivered. Same shape as Shipped, with the handover record
		// in place of the tracking panel. The heading interpolates the name in
		// the template so it stays escaped.
		htmlTemplate = orderDeliveredHTML
		view.Eyebrow = "Delivered"
		view.Intro = "Your parcel was signed for today. Wear it once, then tell us how it sits. Fit adjustments within the first thirty days are always on us, with Oyecouture or through your local tailor with our notes."
	default:
		view.Heading = "Your piece is in production."
		view.Intro = "Work has begun on your Oyecouture order. We will keep you updated as it moves toward dispatch."
	}

	parsedTemplate, err := template.New("fulfillment").Parse(htmlTemplate)
	if err != nil {
		return "", "", err
	}

	var htmlBuffer bytes.Buffer
	if err := parsedTemplate.Execute(&htmlBuffer, view); err != nil {
		return "", "", err
	}

	var textBuilder strings.Builder
	textBuilder.WriteString(view.Heading + "\n\n")
	textBuilder.WriteString("Hi " + view.CustomerName + ",\n")
	textBuilder.WriteString(view.Intro + "\n\n")
	textBuilder.WriteString("Order: " + view.OrderNo + "\n")
	textBuilder.WriteString("Status: " + view.StageLabel + "\n")
	if view.EstimatedDeliveryDate != "" {
		textBuilder.WriteString("Estimated delivery: " + view.EstimatedDeliveryDate + "\n")
	}
	if view.CarrierName != "" {
		textBuilder.WriteString("Courier: " + view.CarrierName + "\n")
	}
	if view.CarrierTrackingNumber != "" {
		textBuilder.WriteString("Courier reference: " + view.CarrierTrackingNumber + "\n")
	}
	if view.PublicNote != "" {
		textBuilder.WriteString("\n" + view.PublicNote + "\n")
	}
	if view.TrackingURL != "" {
		textBuilder.WriteString("\nTrack order: " + view.TrackingURL + "\n")
	}
	if view.CarrierTrackingURL != "" {
		textBuilder.WriteString("Courier tracking: " + view.CarrierTrackingURL + "\n")
	}

	return htmlBuffer.String(), textBuilder.String(), nil
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

	base := strings.TrimRight(siteURL, "/")

	// The customer's own confirmation page, reference included so it can verify
	// rather than reporting a failure to someone who has already paid.
	orderURL := ""
	if order.PayToken != "" && order.Payment.Reference != "" {
		orderURL = base + "/pay/" + order.PayToken + "/callback?reference=" +
			url.QueryEscape(order.Payment.Reference)
	}

	return emailView{
		LogoURL:         base + "/images/brand/oyecouture-logo.png",
		CustomerName:    order.CustomerName,
		FirstName:       firstName(order.CustomerName),
		OrderURL:        orderURL,
		OrderDate:       order.OrderDate,
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
  .action { display: inline-block; padding: 14px 22px; background: #24211f; color: #e0cf9c !important; font-size: 12px; font-weight: 700; letter-spacing: 2.5px; text-decoration: none; text-transform: uppercase; }
  .footer { padding: 26px 38px 34px; color: #8d847c; font-size: 13px; line-height: 1.7; }
  .footer strong { color: #24211f; letter-spacing: 2px; text-transform: uppercase; }
  /* Design 01: dark announcement strip, centred mark, gold rule beneath. */
  .announce { padding: 13px 20px; background: #1c1a18; color: #d9c489; font-size: 11px; font-weight: 700; letter-spacing: 3px; text-align: center; text-transform: uppercase; }
  .brandhead { padding: 30px 38px 26px; background: #faf7f2; border-bottom: 2px solid #a98a3d; text-align: center; }
  .brandtext { color: #24211f; font-size: 26px; font-weight: 700; letter-spacing: 7px; text-transform: uppercase; }
  .brandtext span { display: block; margin-top: 7px; color: #a98a3d; font-size: 11px; font-weight: 600; letter-spacing: 9px; }
  .detailbox { width: 100%; margin-top: 26px; border-collapse: collapse; background: #f7f3ee; border: 1px solid #e8dfd3; }
  .detailbox td { padding: 11px 18px; font-size: 13px; }
  .detailbox tr:first-child td { padding-top: 18px; }
  .detailbox tr:last-child td { padding-bottom: 18px; }
  .dkey { color: #8d847c; letter-spacing: 1.4px; text-transform: uppercase; }
  .dval { color: #24211f; font-weight: 700; text-align: right; }
  .deliverbox { width: 100%; margin-top: 26px; border-collapse: collapse; background: #f7f3ee; border: 1px solid #e8dfd3; }
  .deliverbox td { padding: 6px 18px; font-size: 14px; line-height: 1.6; }
  .deliverbox tr:first-child td { padding-top: 18px; }
  .deliverbox tr:last-child td { padding-bottom: 18px; }
  .signoff { margin-top: 12px; color: #24211f; font-size: 13px; font-weight: 700; letter-spacing: 3px; text-transform: uppercase; }
  .fineprint { color: #9a9188; font-size: 12px; line-height: 1.65; }
  /* Uppercase on the label only: letter-spacing a currency figure reads badly. */
  .totallabel { color: #24211f; font-weight: 700; letter-spacing: 1.6px; text-transform: uppercase; }
  /* Design 02: four-step progress bar. A table, not flex, because Outlook does
     not lay flex out and the steps would stack into a column. */
  .progress { width: 100%; margin-top: 30px; border-collapse: separate; border-spacing: 8px 0; }
  .progress td { padding: 0; vertical-align: top; }
  .progress .bar { display: block; height: 3px; background: #ded5c8; }
  .progress .bar.done { background: #a98a3d; }
  .progress .step { display: block; margin-top: 9px; color: #b3aaa1; font-size: 11px; font-weight: 700; letter-spacing: 1.4px; text-transform: uppercase; }
  .progress .step.done { color: #24211f; }
  @media only screen and (max-width: 600px) {
    .shell { padding: 18px 8px; }
    .header, .section, .footer { padding-left: 22px; padding-right: 22px; }
    h1 { font-size: 25px; }
    .summary td { display: block; width: auto; }
    .announce { font-size: 10px; letter-spacing: 2px; }
    .brandhead { padding-left: 22px; padding-right: 22px; }
  }
`

const fulfillmentUpdateHTML = `<!doctype html>
<html>
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><style>` + emailBaseStyles + `</style></head>
<body>
  <div class="preheader">{{.StageLabel}} for Oyecouture order {{.OrderNo}}.</div>
  <div class="shell">
    <div class="card">
      <div class="topline"></div>
      <div class="header">
        <div class="brand">Oyecouture<span>Bespoke</span></div>
        <div class="eyebrow">Order progress</div>
        <h1>{{.Heading}}</h1>
        <table class="summary">
          <tr>
            <td>Order<strong>{{.OrderNo}}</strong></td>
            <td>Status<br><span class="badge">{{.StageLabel}}</span></td>
          </tr>
        </table>
      </div>
      <div class="section">
        <p>Hi {{.CustomerName}}, {{.Intro}}</p>
      </div>
      {{if or .EstimatedDeliveryDate .CarrierName .CarrierTrackingNumber}}
      <div class="section alt">
        <div class="label">Delivery details</div>
        <table class="row">
          {{if .EstimatedDeliveryDate}}<tr><td>Estimated delivery</td><td class="value">{{.EstimatedDeliveryDate}}</td></tr>{{end}}
          {{if .CarrierName}}<tr><td>Courier</td><td class="value">{{.CarrierName}}</td></tr>{{end}}
          {{if .CarrierTrackingNumber}}<tr><td>Courier reference</td><td class="value">{{.CarrierTrackingNumber}}</td></tr>{{end}}
        </table>
      </div>
      {{end}}
      {{if .PublicNote}}
      <div class="section"><div class="note">{{.PublicNote}}</div></div>
      {{end}}
      {{if .TrackingURL}}
      <div class="section"><a class="action" href="{{.TrackingURL}}">Track your order</a></div>
      {{end}}
      {{if .CarrierTrackingURL}}
      <div class="section alt"><p>You can also <a href="{{.CarrierTrackingURL}}">view updates from the courier</a>.</p></div>
      {{end}}
      <div class="footer"><strong>Oyecouture Bespoke</strong><br>This update was sent from your Oyecouture order timeline.</div>
    </div>
  </div>
</body>
</html>`

// orderConfirmationHTML follows the client's approved "01 Payment received"
// design. Table-based and inline-friendly because email clients are not
// browsers: flex and grid are unreliable in Outlook, so structure comes from
// tables and the shared styles above.
const orderConfirmationHTML = `<!doctype html>
<html>
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><style>` + emailBaseStyles + `</style></head>
<body>
  <div class="preheader">We have your payment and your pieces are entering the workroom.</div>
  <div class="shell">
    <div class="card">
      <div class="announce">Complimentary tailoring consultation &middot; Worldwide delivery</div>

      <!-- Set as type, not an image. Mail clients block remote images by
           default, so a hosted logo means most recipients open an email with
           no brand on it at all. -->
      <div class="brandhead">
        <div class="brandtext">Oyecouture<span>Bespoke</span></div>
      </div>

      <div class="section">
        <p class="eyebrow">Payment received</p>
        <h1>Thank you, {{.FirstName}}.<br>Your order is confirmed.</h1>
        <p style="margin-top:18px">Your payment has come through and your commission is now with our cutters. We will write again the moment your pieces leave us, and if anything needs clarifying before then, simply reply to this email.</p>

        <table class="detailbox" role="presentation">
          <tr><td colspan="2" class="label" style="padding-bottom:4px">Order details</td></tr>
          <tr><td class="dkey">Order reference</td><td class="dval">{{.OrderNo}}</td></tr>
          <tr><td class="dkey">Order date</td><td class="dval">{{.OrderDate}}</td></tr>
          <tr><td class="dkey">Payment</td><td class="dval">{{.PaymentProvider}} &middot; Verified</td></tr>
        </table>
      </div>

      <div class="section">
        <p class="label">Your pieces</p>
        <table class="items" role="presentation">
          {{range .Items}}
          <tr>
            <td>{{.Name}}<span class="muted">{{if .Size}}Size {{.Size}} &middot; {{end}}Qty {{.Quantity}}</span></td>
            <td class="money">{{.Line}}</td>
          </tr>
          {{end}}
        </table>

        <table class="row" role="presentation" style="margin-top:18px">
          <tr><td>Subtotal</td><td class="value">{{.Subtotal}}</td></tr>
          <tr><td>Tailoring consultation</td><td class="value">Complimentary</td></tr>
          <tr><td>Delivery</td><td class="value">{{.DeliveryFee}}</td></tr>
          <tr class="grand-total"><td class="total totallabel">Total paid</td><td class="value total">{{.Total}}</td></tr>
        </table>

        <table class="deliverbox" role="presentation">
          <tr><td class="label" style="padding-bottom:10px">Delivering to</td></tr>
          <tr><td style="color:#24211f">{{.CustomerName}}<br>{{.Address}}{{if .DeliveryZone}}<br>{{.DeliveryZone}}{{end}}{{if .Phone}}<br>{{.Phone}}{{end}}</td></tr>
        </table>

        {{if .OrderURL}}
        <div style="text-align:center;margin-top:30px">
          <a class="action" href="{{.OrderURL}}">View your order</a>
        </div>
        {{end}}
      </div>

      <div class="section" style="border-bottom:0">
        <p>We are glad to be dressing you.</p>
        <p class="signoff">Oyecouture Bespoke</p>
      </div>

      <div class="footer">
        <p style="margin-bottom:14px">Questions about this order? Reply to this email and a stylist will answer personally.</p>
        <p style="margin-bottom:14px">Ready-to-wear pieces may be exchanged within 7 days of delivery in original condition. Made-to-measure and bespoke commissions are cut to your measurements and are final sale, though we will adjust the fit at no charge.</p>
        <p class="fineprint">Oyecouture Bespoke &middot; 6 Force Road, Tafawa Balewa Square (TBS), Lagos, Nigeria<br>You received this email because you placed order {{.OrderNo}}.</p>
      </div>
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

// orderDispatchedHTML follows the client's approved "02 Shipped" design: the
// four-step progress bar, a tracking panel, the order reproduced, and a link to
// follow the parcel.
const orderDispatchedHTML = `<!doctype html>
<html>
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><style>` + emailBaseStyles + `</style></head>
<body>
  <div class="preheader">Your pieces left us today. Tracking details inside.</div>
  <div class="shell">
    <div class="card">
      <div class="announce">Complimentary tailoring consultation &middot; Worldwide delivery</div>

      <div class="brandhead">
        <div class="brandtext">Oyecouture<span>Bespoke</span></div>
      </div>

      <div class="section">
        <p class="eyebrow">{{.Eyebrow}}</p>
        <h1>{{.FirstName}}, your order<br>has left us.</h1>
        <p style="margin-top:18px">{{.Intro}}</p>

        <table class="progress" role="presentation">
          <tr>
            {{range .Progress}}<td width="25%"><span class="bar{{if .Done}} done{{end}}"></span></td>{{end}}
          </tr>
          <tr>
            {{range .Progress}}<td width="25%"><span class="step{{if .Done}} done{{end}}">{{.Label}}</span></td>{{end}}
          </tr>
        </table>

        {{if or .CarrierName .CarrierTrackingNumber .EstimatedDeliveryDate}}
        <table class="detailbox" role="presentation">
          <tr><td colspan="2" class="label" style="padding-bottom:4px">Tracking</td></tr>
          {{if .CarrierName}}<tr><td class="dkey">Courier</td><td class="dval">{{.CarrierName}}</td></tr>{{end}}
          {{if .CarrierTrackingNumber}}<tr><td class="dkey">Tracking number</td><td class="dval">{{.CarrierTrackingNumber}}</td></tr>{{end}}
          {{if .EstimatedDeliveryDate}}<tr><td class="dkey">Estimated arrival</td><td class="dval">{{.EstimatedDeliveryDate}}</td></tr>{{end}}
        </table>
        {{end}}

        {{if .PublicNote}}<div class="note" style="margin-top:24px">{{.PublicNote}}</div>{{end}}
      </div>

      {{if .Items}}
      <div class="section">
        <p class="label">Your pieces</p>
        <table class="items" role="presentation">
          {{range .Items}}
          <tr>
            <td>{{.Name}}<span class="muted">{{if .Size}}Size {{.Size}} &middot; {{end}}Qty {{.Quantity}}</span></td>
            <td class="money">{{.Line}}</td>
          </tr>
          {{end}}
        </table>

        <table class="row" role="presentation" style="margin-top:18px">
          <tr><td>Subtotal</td><td class="value">{{.Subtotal}}</td></tr>
          <tr><td>Tailoring consultation</td><td class="value">Complimentary</td></tr>
          <tr><td>Delivery</td><td class="value">{{.DeliveryFee}}</td></tr>
          <tr class="grand-total"><td class="total totallabel">Total paid</td><td class="value total">{{.Total}}</td></tr>
        </table>
      </div>
      {{end}}

      <!-- Outside the items block on purpose. Following the parcel is the whole
           point of this email, so the link must not depend on the order
           reproducing. Our own tracking page is preferred over the courier's:
           it always exists, and it is the link we control. -->
      {{if or .TrackingURL .CarrierTrackingURL}}
      <div class="section" style="text-align:center">
        <a class="action" href="{{if .TrackingURL}}{{.TrackingURL}}{{else}}{{.CarrierTrackingURL}}{{end}}">Track your parcel</a>
      </div>
      {{end}}

      <div class="section" style="border-bottom:0">
        <p>Thank you for your patience while it was made.</p>
        <p class="signoff">Oyecouture Bespoke</p>
      </div>

      <div class="footer">
        <p style="margin-bottom:14px">Questions about this order? Reply to this email and a stylist will answer personally.</p>
        <p style="margin-bottom:14px">Ready-to-wear pieces may be exchanged within 7 days of delivery in original condition. Made-to-measure and bespoke commissions are cut to your measurements and are final sale, though we will adjust the fit at no charge.</p>
        <p class="fineprint">Oyecouture Bespoke &middot; 6 Force Road, Tafawa Balewa Square (TBS), Lagos, Nigeria<br>You received this email because you placed order {{.OrderNo}}.</p>
      </div>
    </div>
  </div>
</body>
</html>`

// orderDeliveredHTML follows the client's approved "04 Delivered" design: the
// completed progress bar, the handover record, and the order reproduced. The
// design's "Leave a review" panel is deliberately absent - there is no review
// system for it to reach, and a button that goes nowhere is worse than none.
const orderDeliveredHTML = `<!doctype html>
<html>
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><style>` + emailBaseStyles + `</style></head>
<body>
  <div class="preheader">Your pieces have arrived. A note on fit adjustments.</div>
  <div class="shell">
    <div class="card">
      <div class="announce">Complimentary tailoring consultation &middot; Worldwide delivery</div>

      <div class="brandhead">
        <div class="brandtext">Oyecouture<span>Bespoke</span></div>
      </div>

      <div class="section">
        <p class="eyebrow">{{.Eyebrow}}</p>
        <h1>{{.FirstName}}, your order<br>has arrived.</h1>
        <p style="margin-top:18px">{{.Intro}}</p>

        <table class="progress" role="presentation">
          <tr>
            {{range .Progress}}<td width="25%"><span class="bar{{if .Done}} done{{end}}"></span></td>{{end}}
          </tr>
          <tr>
            {{range .Progress}}<td width="25%"><span class="step{{if .Done}} done{{end}}">{{.Label}}</span></td>{{end}}
          </tr>
        </table>

        {{if or .SignedForBy .SignedAt .DeliveryLocation}}
        <table class="detailbox" role="presentation">
          <tr><td colspan="2" class="label" style="padding-bottom:4px">Delivered</td></tr>
          {{if .SignedForBy}}<tr><td class="dkey">Signed for by</td><td class="dval">{{.SignedForBy}}</td></tr>{{end}}
          {{if .SignedAt}}<tr><td class="dkey">Date</td><td class="dval">{{.SignedAt}}</td></tr>{{end}}
          {{if .DeliveryLocation}}<tr><td class="dkey">Location</td><td class="dval">{{.DeliveryLocation}}</td></tr>{{end}}
        </table>
        {{end}}

        {{if .PublicNote}}<p class="note" style="margin-top:26px">{{.PublicNote}}</p>{{end}}
      </div>

      {{if .Items}}
      <div class="section">
        <p class="label">Your pieces</p>
        <table class="items" role="presentation">
          {{range .Items}}
          <tr>
            <td>{{.Name}}<span class="muted">{{if .Size}}Size {{.Size}} &middot; {{end}}Qty {{.Quantity}}</span></td>
            <td class="money">{{.Line}}</td>
          </tr>
          {{end}}
        </table>

        <table class="row" role="presentation" style="margin-top:18px">
          <tr><td>Subtotal</td><td class="value">{{.Subtotal}}</td></tr>
          <tr><td>Tailoring consultation</td><td class="value">Complimentary</td></tr>
          <tr><td>Delivery</td><td class="value">{{.DeliveryFee}}</td></tr>
          <tr class="grand-total"><td class="total totallabel">Total paid</td><td class="value total">{{.Total}}</td></tr>
        </table>
      </div>
      {{end}}

      <div class="section" style="border-bottom:0">
        <p>Thank you for wearing Oyecouture.</p>
        <p class="signoff">Oyecouture Bespoke</p>
      </div>

      <div class="footer">
        <p style="margin-bottom:14px">Questions about this order? Reply to this email and a stylist will answer personally.</p>
        <p style="margin-bottom:14px">Ready-to-wear pieces may be exchanged within 7 days of delivery in original condition. Made-to-measure and bespoke commissions are cut to your measurements and are final sale, though we will adjust the fit at no charge.</p>
        <p class="fineprint">Oyecouture Bespoke &middot; 6 Force Road, Tafawa Balewa Square (TBS), Lagos, Nigeria<br>You received this email because you placed order {{.OrderNo}}.</p>
      </div>
    </div>
  </div>
</body>
</html>`

// outForDeliveryHTML follows the client's approved "03 Out for delivery"
// design: the rider carrying the parcel, their number, and the window to expect
// them. Distinct from dispatch, which names a courier and a tracking number;
// this step is Oyecouture's own rider on a Lagos road, and the customer needs a
// person to expect rather than a parcel to track.
const outForDeliveryHTML = `<!doctype html>
<html>
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><style>` + emailBaseStyles + `</style></head>
<body>
  <div class="preheader">Your parcel is with the rider and arrives today.</div>
  <div class="shell">
    <div class="card">
      <div class="announce">Complimentary tailoring consultation &middot; Worldwide delivery</div>

      <div class="brandhead">
        <div class="brandtext">Oyecouture<span>Bespoke</span></div>
      </div>

      <div class="section">
        <p class="eyebrow">{{.Eyebrow}}</p>
        <h1>{{.FirstName}}, your parcel<br>is out for delivery.</h1>
        <p style="margin-top:18px">{{.Intro}}</p>

        <table class="progress" role="presentation">
          <tr>
            {{range .Progress}}<td width="25%"><span class="bar{{if .Done}} done{{end}}"></span></td>{{end}}
          </tr>
          <tr>
            {{range .Progress}}<td width="25%"><span class="step{{if .Done}} done{{end}}">{{.Label}}</span></td>{{end}}
          </tr>
        </table>

        {{if or .RiderName .RiderContact .DeliveryWindow}}
        <table class="detailbox" role="presentation">
          <tr><td colspan="2" class="label" style="padding-bottom:4px">Today's delivery</td></tr>
          {{if .RiderName}}<tr><td class="dkey">Rider</td><td class="dval">{{.RiderName}}</td></tr>{{end}}
          {{if .RiderContact}}<tr><td class="dkey">Contact</td><td class="dval">{{.RiderContact}}</td></tr>{{end}}
          {{if .DeliveryWindow}}<tr><td class="dkey">Window</td><td class="dval">{{.DeliveryWindow}}</td></tr>{{end}}
        </table>
        {{end}}

        {{if .PublicNote}}<p class="note" style="margin-top:26px">{{.PublicNote}}</p>{{end}}

        {{if .TrackingURL}}
        <div style="text-align:center;margin-top:30px">
          <a class="action" href="{{.TrackingURL}}">Follow the delivery</a>
        </div>
        {{end}}
      </div>

      <div class="section" style="border-bottom:0">
        <p>We hope it is everything you pictured.</p>
        <p class="signoff">Oyecouture Bespoke</p>
      </div>

      <div class="footer">
        <p style="margin-bottom:14px">Questions about this order? Reply to this email and a stylist will answer personally.</p>
        <p class="fineprint">Oyecouture Bespoke &middot; 6 Force Road, Tafawa Balewa Square (TBS), Lagos, Nigeria<br>You received this email because you placed order {{.OrderNo}}.</p>
      </div>
    </div>
  </div>
</body>
</html>`
