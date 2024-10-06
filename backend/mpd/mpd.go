package mpd

import (
	"encoding/xml"
	"io"
	"os"
)

type MPD struct {
	XMLName                   xml.Name        `xml:"urn:mpeg:DASH:schema:MPD:2011 MPD"`
	SchemaLocation            string          `xml:"http://www.w3.org/2001/XMLSchema-instance schemaLocation,attr"`
	Type                      string          `xml:"type,attr"`
	MediaPresentationDuration ISO8601Duration `xml:"mediaPresentationDuration,attr"`
	MinBufferTime             ISO8601Duration `xml:"minBufferTime,attr"`
	Profiles                  string          `xml:"profiles,attr"`
	Period                    []Period
}

type Period struct {
	Id            int             `xml:"id,attr"`
	Start         ISO8601Duration `xml:"start,attr"`
	Duration      ISO8601Duration `xml:"duration,attr"`
	AdaptationSet []AdaptationSet
}

type AdaptationSet struct {
	Id                      int    `xml:"id,attr"`
	MimeType                string `xml:"mimeType,attr"`
	Codecs                  string `xml:"codecs,attr"`
	Lang                    string `xml:"lang,attr"`
	Width                   int    `xml:"width,attr,omitempty"`
	Height                  int    `xml:"height,attr,omitempty"`
	BitstreamSwitching      bool   `xml:"bitstreamSwitching,attr"`
	SubsegmentAlignment     bool   `xml:"subsegmentAlignment,attr"`
	SubsegmentStartsWithSAP int    `xml:"subsegmentStartsWithSAP,attr"`
	Representation          []Representation
}

type Representation struct {
	Id          int `xml:"id,attr"`
	Bandwidth   int `xml:"bandwidth,attr"`
	Width       int `xml:"width,attr,omitempty"`
	Height      int `xml:"height,attr,omitempty"`
	BaseURL     string
	SegmentBase SegmentBase
}

type SegmentBase struct {
	IndexRange     string `xml:"indexRange,attr"`
	Initialization Initialization
}

type Initialization struct {
	Range string `xml:"range,attr"`
}

func Decode(r io.Reader) (MPD, error) {
	var mpd MPD
	err := xml.NewDecoder(r).Decode(&mpd)
	return mpd, err
}

func FromFile(name string) (MPD, error) {
	f, err := os.Open(name)
	if err != nil {
		return MPD{}, err
	}
	defer f.Close()
	return Decode(f)
}

func (m MPD) Encode(w io.Writer) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(m)
}

func (m MPD) WriteFile(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return m.Encode(f)
}
