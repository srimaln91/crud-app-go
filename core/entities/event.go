package entities

type Event struct {
	ID          string `json:"eventId"`
	TransId     string `json:"transId"`
	TransTms    string `json:"transTms"`
	RcNum       string `json:"rcNum"`
	ClientId    string `json:"clientId"`
	EventCnt    int    `json:"eventCnt"`
	LocationCd  string `json:"locationCd"`
	LocationId1 string `json:"locationId1"`
	LocationId2 string `json:"locationId2"`
	AddrNbr     string `json:"addrNbr"`
}
