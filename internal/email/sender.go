// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package email

import (
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/superseriousbusiness/gotosocial/internal/config"
)

// Sender contains functions for sending emails to instance users/new signups.
type Sender interface {
	// SendConfirmEmail sends a 'please confirm your email' style email to the given toAddress, with the given data.
	SendConfirmEmail(toAddress string, data ConfirmData) error

	// SendResetEmail sends a 'reset your password' style email to the given toAddress, with the given data.
	SendResetEmail(toAddress string, data ResetData) error

	// SendTestEmail sends a 'testing email sending' style email to the given toAddress, with the given data.
	SendTestEmail(toAddress string, data TestData) error

	// SendNewReportEmail sends an email notification to the given addresses, letting them
	// know that a new report has been created targeting a user on this instance.
	//
	// It is expected that the toAddresses have already been filtered to ensure that they
	// all belong to admins + moderators.
	SendNewReportEmail(toAddresses []string, data NewReportData) error

	// SendReportClosedEmail sends an email notification to the given address, letting them
	// know that a report that they created has been closed / resolved by an admin.
	SendReportClosedEmail(toAddress string, data ReportClosedData) error
}

// NewSender returns a new email Sender interface with the given configuration, or an error if something goes wrong.
func NewSender() (Sender, error) {
	templateBaseDir := config.GetWebTemplateBaseDir()
	t, err := loadTemplates(templateBaseDir)
	if err != nil {
		return nil, err
	}

	username := config.GetSMTPUsername()
	password := config.GetSMTPPassword()
	host := config.GetSMTPHost()
	port := config.GetSMTPPort()
	from := config.GetSMTPFrom()

	return &sender{
		hostAddress: fmt.Sprintf("%s:%d", host, port),
		from:        from,
		auth:        smtp.PlainAuth("", username, password, host),
		template:    t,
	}, nil
}

type sender struct {
	hostAddress string
	from        string
	auth        smtp.Auth
	template    *template.Template
}
