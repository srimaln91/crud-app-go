package request

type EventBatch struct {
	BatchID string `json:"batchId"`
	Records []struct {
		TransID  string `json:"transId"`
		TransTms string `json:"transTms"`
		RcNum    string `json:"rcNum"`
		ClientID string `json:"clientId"`
		Event    []struct {
			EventCnt    int    `json:"eventCnt"`
			LocationCd  string `json:"locationCd"`
			LocationID1 string `json:"locationId1"`
			LocationID2 string `json:"locationId2,omitempty"`
			AddrNbr     string `json:"addrNbr,omitempty"`
		} `json:"event"`
	} `json:"records"`
}
