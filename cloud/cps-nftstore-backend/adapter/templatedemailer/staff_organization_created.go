package templatedemailer

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"text/template"

	"log/slog"
)

func (impl *templatedEmailer) SendNewStoreEmailToStaff(staffEmails []string, storeID string) error {
	impl.Logger.Debug("sending `New Registration` to admin staff", slog.String("storeID", storeID))

	for _, staffEmail := range staffEmails {
		// FOR TESTING PURPOSES ONLY.
		fp := path.Join("templates", "staff_store_created.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			impl.Logger.Error("parsing error", slog.Any("error", err))
			return err
		}

		var processed bytes.Buffer

		// Render the HTML template with our data.
		data := struct {
			DetailLink string
		}{
			DetailLink: fmt.Sprintf("https://%v/admin/store/%v", impl.Emailer.GetDomainName(), storeID),
		}
		if err := tmpl.Execute(&processed, data); err != nil {
			impl.Logger.Error("template execution error", slog.Any("error", err))
			return err
		}
		body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

		if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "New CPS Retail Partner Application", staffEmail, body); err != nil {
			impl.Logger.Error("sending error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("sent `New Registration` email",
			slog.String("staffEmail", staffEmail),
			slog.Any("storeID", storeID))
	}
	return nil
}
