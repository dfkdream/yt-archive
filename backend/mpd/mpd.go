package mpd

import (
	"encoding/xml"
	"io"
	"os"
)

type MPD struct {
	XMLName                   xml.Name `xml:"urn:mpeg:DASH:schema:MPD:2011 MPD"`
	SchemaLocation            string   `xml:"http://www.w3.org/2001/XMLSchema-instance schemaLocation,attr"`
	Type                      string   `xml:"type,attr"`
	MediaPresentationDuration string   `xml:"mediaPresentationDuration,attr"`
	MinBufferTime             string   `xml:"minBufferTime,attr"`
	Profiles                  string   `xml:"profiles,attr"`
	Period                    []Period
}

type Period struct {
	Id            string `xml:"id,attr"`
	Start         string `xml:"start,attr"`
	Duration      string `xml:"duration,attr"`
	AdaptationSet []AdaptationSet
}

type AdaptationSet struct {
	Id                      string `xml:"id,attr"`
	MimeType                string `xml:"mimeType,attr"`
	Codecs                  string `xml:"codecs,attr"`
	Lang                    string `xml:"lang,attr"`
	Width                   int    `xml:"width,attr,omitempty"`
	Height                  int    `xml:"height,attr,omitempty"`
	BitstreamSwitching      bool   `xml:"bitstreamSwitching,attr"`
	SubsegmentAlignment     bool   `xml:"subsegmentAlignment,attr"`
	SubsegmentStartsWithSAP string `xml:"subsegmentStartsWithSAP,attr"`
	Representation          []Representation
}

type Representation struct {
	Id          string `xml:"id,attr"`
	Bandwidth   string `xml:"bandwidth,attr"`
	Width       int    `xml:"width,attr,omitempty"`
	Height      int    `xml:"height,attr,omitempty"`
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
