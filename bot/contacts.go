package bot

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type User struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type TrueUser struct {
	Name     string
	Position string
	Company  string
	Mail     string
}

type DirectoryData struct {
	XMLName xml.Name `xml:"directoryData"`
	Users   []User   `xml:"field"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

type XMLBODY struct {
	XMLName       xml.Name        `xml:"XIMSS"`
	DirectoryData []DirectoryData `xml:"directoryData"`
}

const xmlbody = `
<XMLBODY>
	<listDirectory limit="%s" filter="(|(mail=*%s*))" id="%s">
		<field>cn</field>
		<field>mail</field>
		<field>o</field>
		<field>ou</field>
		<field>telephoneNumber</field>
		<field>mobile</field>
		<field>userCertificate</field>
	</listDirectory>
</XMLBODY>`

func TimeStamp() string {
	ts := time.Now().Unix()
	return fmt.Sprint(ts)
}

func BuildMessageFromUsers(users []TrueUser) string {
	buff := ""
	for _, user := range users {
		buff += fmt.Sprintf(
			"---\nИмя: %s\nДолжность: %s\nКомпания: %s\nПочта: %s\n",
			user.Name, user.Position, user.Company, user.Mail,
		)
	}
	return buff
}

func buildXML(mail, count string) string {
	return fmt.Sprintf(xmlbody, count, mail, TimeStamp())
}

func GetUserByRegex(regex, url, count string) ([]TrueUser, error) {
	resp, err := http.Post(url, "text/xml", strings.NewReader(buildXML(regex, count)))
	if err != nil {
		return []TrueUser{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []TrueUser{}, err
	}

	var usersDry XMLBODY
	var users []TrueUser
	xml.Unmarshal(body, &usersDry)

	for _, dir := range usersDry.DirectoryData {
		var tmp []string
		for _, user := range dir.Users {
			tmp = append(tmp, user.Value)
		}
		users = append(
			users,
			TrueUser{
				Name:     tmp[0],
				Position: tmp[1],
				Company:  tmp[2],
				Mail:     tmp[3],
			},
		)
	}

	return users, nil
}
