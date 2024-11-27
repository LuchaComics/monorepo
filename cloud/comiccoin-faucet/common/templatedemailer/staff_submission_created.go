package templatedemailer

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"text/template"

	"log/slog"
)

func (impl *templatedEmailer) SendNewComicSubmissionEmailToStaff(staffEmails []string, submissionID string, storeName string, item string, cpsrn string, serviceTypeName string) error {
	impl.Logger.Debug("sending `New Comic Submission` to admin staff", slog.String("submissionID", submissionID))

	for _, staffEmail := range staffEmails {
		// FOR TESTING PURPOSES ONLY.
		fp := path.Join("templates", "staff_submission_created.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			impl.Logger.Error("parsing error",
				slog.Any("error", err))
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
			DetailLink:      fmt.Sprintf("https://%v/admin/submission/%v", impl.Emailer.GetDomainName(), submissionID),
		}
		if err := tmpl.Execute(&processed, data); err != nil {
			impl.Logger.Error("template execution error",
				slog.Any("data", data),
				slog.String("StoreName", data.StoreName),
				slog.String("Item", data.Item),
				slog.String("CPSRN", data.CPSRN),
				slog.String("DetailLink", data.DetailLink),
				slog.Any("error", err),
			)
			return err
		}
		body := processed.String() // DEVELOPERS NOTE: Convert our long sequence of data into a string.

		if err := impl.Emailer.Send(context.Background(), impl.Emailer.GetSenderEmail(), "New Comic Submission", staffEmail, body); err != nil {
			impl.Logger.Error("sending error",
				slog.Any("staffEmail", staffEmail),
				slog.Any("submissionID", submissionID),
				slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("sent `New Comic Submission` email",
			slog.String("staffEmail", staffEmail),
			slog.Any("submissionID", submissionID))
	}
	return nil
}
