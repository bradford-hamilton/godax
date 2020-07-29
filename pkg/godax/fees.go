package godax

// Fees represent the current maker & taker fees, and usd volume.
/*
{
    "maker_fee_rate": "0.0015",
    "taker_fee_rate": "0.0025",
    "usd_volume": "25000.00"
}
*/
type Fees struct {
	MakerFeeRate string `json:"maker_fee_rate"`
	TakerFeeRate string `json:"taker_fee_rate"`
	USDVolume    string `json:"usd_volume"`
}
