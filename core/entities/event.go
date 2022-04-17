package entities

type Event struct {
	ID          string `json:"eventId"`
	TransId     string `json:"transId,omitempty"`
	TransTms    string `json:"transTms,omitempty"`
	RcNum       string `json:"rcNum,omitempty"`
	ClientId    string `json:"clientId,omitempty"`
	EventCnt    int    `json:"eventCnt,omitempty"`
	LocationCd  string `json:"locationCd,omitempty"`
	LocationId1 string `json:"locationId1,omitempty"`
	LocationId2 string `json:"locationId2,omitempty"`
	AddrNbr     string `json:"addrNbr,omitempty"`
}
