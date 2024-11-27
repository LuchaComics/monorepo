package templatedemailer

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"text/template"

	"log/slog"
)

func (impl *templatedEmailer) SendNewComicSubmissionEmailToRetailers(retailerEmails []string, submissionID string, storeName string, item string, cpsrn string, serviceTypeName string) error {
	impl.Logger.Debug("sending `Submitted to CPS` to retailer", slog.String("submissionID", submissionID))

	for _, retailerEmail := range retailerEmails {
		// FOR TESTING PURPOSES ONLY.
		fp := path.Join("templates", "retailer_submission_created.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			impl.Logger.Error("parsing error", slog.Any("error", err))
			return err
		}

		var processed bytes.Buffer

		// Render the HTML template with our data.
		data := struct {
			StoreName       string
			Item            string
			CPSRN           string
			ServiceTypeName string
			DetailLink      string
		}{
			StoreName:       storeName,
			Item:            item,
			CPSRN:           cpsrn,
			ServiceTypeName: serviceTypeName,
			DetailLink:      fmt.Sprintf("https://%v/submission/%v", impl.Emailer.GetDomainName(), submissionID),
		}
		if err := tmpl.Execute(&processed, data); err != nil {
			impl.Logger.Error("template execution error", slog.Any("error", err))
			return err
		}
		body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

		if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "New Comic Submission", retailerEmail, body); err != nil {
			impl.Logger.Error("sending error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("sent `Submitted to CPS` email to retailer",
			slog.String("retailerEmail", retailerEmail),
			slog.Any("submissionID", submissionID))
	}
	return nil
}
