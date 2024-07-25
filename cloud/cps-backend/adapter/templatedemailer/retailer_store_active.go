package templatedemailer

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"text/template"

	"log/slog"
)

func (impl *templatedEmailer) SendRetailerStoreActiveEmailToRetailers(retailerEmails []string, storeName string) error {
	impl.Logger.Debug("sending `Store Active` email to retailer")

	for _, retailerEmail := range retailerEmails {
		// FOR TESTING PURPOSES ONLY.
		fp := path.Join("templates", "retailer_store_active.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			impl.Logger.Error("parsing error", slog.Any("error", err))
			return err
		}

		var processed bytes.Buffer

		// Render the HTML template with our data.
		data := struct {
			StoreName  string
			DetailLink string
		}{
			StoreName:  storeName,
			DetailLink: fmt.Sprintf("https://%v/login", impl.Emailer.GetDomainName()),
		}
		if err := tmpl.Execute(&processed, data); err != nil {
			impl.Logger.Error("template execution error", slog.Any("error", err))
			return err
		}
		body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

		if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "Your store is active", retailerEmail, body); err != nil {
			impl.Logger.Error("sending error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("sent `Store Active` email  to retailer",
			slog.String("retailerEmail", retailerEmail))
	}
	return nil
}
