package bot

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type User struct {
	Name     string
	Position string
	Company  string
	Mail     string
}

const xmlbody = `
	<XIMSS>
		<listDirectory limit="4" filter="(|(mail=%s*))" id="%s">
			<field>cn</field>
			<field>mail</field>
			<field>o</field>
			<field>ou</field>
			<field>telephoneNumber</field>
			<field>mobile</field>
			<field>userCertificate</field>
		</listDirectory>
	</XIMSS>`

func buildXML(mail string) string {
	return fmt.Sprintf(xmlbody, mail, time.Stamp)
}

func GetUserByMail(mail, url string) (User, error) {
	resp, err := http.Post(url, "text/xml", strings.NewReader(buildXML(mail)))
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

}
