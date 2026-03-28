package model

type Message struct {
	MsgId          string `parquet:"MsgId"`
	PartitionId    uint64 `parquet:"PartitionId"`
	Timestamp      string `parquet:"Timestamp"`
	Hostname       string `parquet:"Hostname"`
	Priority       int32  `parquet:"Priority"`
	Facility       int32  `parquet:"Facility"`
	FacilityString string `parquet:"FacilityString"`
	Severity       int32  `parquet:"Severity"`
	SeverityString string `parquet:"SeverityString"`
	AppName        string `parquet:"AppName"`
	ProcId         string `parquet:"ProcId"`
	Message        string `parquet:"Message"`
	MessageRaw     string `parquet:"MessageRaw"`
	StructuredData string `parquet:"StructuredData"`
	Tag            string `parquet:"Tag"`
	Sender         string `parquet:"Sender"`
	Groupings      string `parquet:"Groupings"`
	Event          string `parquet:"Event"`
	EventId        string `parquet:"EventId"`
	NanoTimeStamp  string `parquet:"NanoTimeStamp"`
	Namespace      string `parquet:"namespace"`
}
