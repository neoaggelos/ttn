// Copyright © 2015 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

// Package gateway/protocol provides useful methods and types to handle communications with a gateway.
//
// This package relies on the SemTech Protocol 1.2 accessible on github: https://github.com/TheThingsNetwork/packet_forwarder/blob/master/PROTOCOL.TXT
package protocol

import (
	"time"
)

// RXPK represents an uplink json message format sent by the gateway
type RXPK struct {
	Chan uint      `json:"chan"` // Concentrator "IF" channel used for RX (unsigned integer)
	Codr string    `json:"codr"` // LoRa ECC coding rate identifier
	Data string    `json:"data"` // Base64 encoded RF packet payload, padded
	Datr string    `json:"-"`    // FSK datarate (unsigned in bit per second) || LoRa datarate identifier
	Freq float64   `json:"freq"` // RX Central frequency in MHx (unsigned float, Hz precision)
	Lsnr float64   `json:"lsnr"` // LoRa SNR ratio in dB (signed float, 0.1 dB precision)
	Modu string    `json:"modu"` // Modulation identifier "LORA" or "FSK"
	Rfch uint      `json:"rfch"` // Concentrator "RF chain" used for RX (unsigned integer)
	Rssi int       `json:"rssi"` // RSSI in dBm (signed integer, 1 dB precision)
	Size uint      `json:"size"` // RF packet payload size in bytes (unsigned integer)
	Stat int       `json:"stat"` // CRC status: 1 - OK, -1 = fail, 0 = no CRC
	Time time.Time `json:"-"`    // UTC time of pkt RX, us precision, ISO 8601 'compact' format
	Tmst uint      `json:"tmst"` // Internal timestamp of "RX finished" event (32b unsigned)
}

// TXPK represents a downlink json message format received by the gateway.
// Most field are optional.
type TXPK struct {
	Codr string    `json:"codr"` // LoRa ECC coding rate identifier
	Data string    `json:"data"` // Base64 encoded RF packet payload, padding optional
	Datr string    `json:"-"`    // LoRa datarate identifier (eg. SF12BW500) || FSK Datarate (unsigned, in bits per second)
	Fdev uint      `json:"fdev"` // FSK frequency deviation (unsigned integer, in Hz)
	Freq float64   `json:"freq"` // TX central frequency in MHz (unsigned float, Hz precision)
	Imme bool      `json:"imme"` // Send packet immediately (will ignore tmst & time)
	Ipol bool      `json:"ipol"` // Lora modulation polarization inversion
	Modu string    `json:"modu"` // Modulation identifier "LORA" or "FSK"
	Ncrc bool      `json:"ncrc"` // If true, disable the CRC of the physical layer (optional)
	Powe uint      `json:"powe"` // TX output power in dBm (unsigned integer, dBm precision)
	Prea uint      `json:"prea"` // RF preamble size (unsigned integer)
	Rfch uint      `json:"rfch"` // Concentrator "RF chain" used for TX (unsigned integer)
	Size uint      `json:"size"` // RF packet payload size in bytes (unsigned integer)
	Time time.Time `json:"-"`    // Send packet at a certain time (GPS synchronization required)
	Tmst uint      `json:"tmst"` // Send packet on a certain timestamp value (will ignore time)
}

// Stat represents a status json message format sent by the gateway
type Stat struct {
	Ackr float64   `json:"ackr"` // Percentage of upstream datagrams that were acknowledged
	Alti int       `json:"alti"` // GPS altitude of the gateway in meter RX (integer)
	Dwnb uint      `json:"dwnb"` // Number of downlink datagrams received (unsigned integer)
	Lati float64   `json:"lati"` // GPS latitude of the gateway in degree (float, N is +)
	Long float64   `json:"long"` // GPS latitude of the gateway in dgree (float, E is +)
	Rxfw uint      `json:"rxfw"` // Number of radio packets forwarded (unsigned integer)
	Rxnb uint      `json:"rxnb"` // Number of radio packets received (unsigned integer)
	Rxok uint      `json:"rxok"` // Number of radio packets received with a valid PHY CRC
	Time time.Time `json:"-"`    // UTC 'system' time of the gateway, ISO 8601 'expanded' format
	Txnb uint      `json:"txnb"` // Number of packets emitted (unsigned integer)
}

// Packet as seen by the gateway.
type Packet struct {
	Version    byte     // Protocol version, should always be 1 here
	Token      []byte   // Random number generated by the gateway on some request. 2-bytes long.
	Identifier byte     // Packet's command identifier
	GatewayId  []byte   // Source gateway's identifier (Only PULL_DATA and PUSH_DATA)
	Payload    *Payload // JSON payload transmitted if any, nil otherwise
}

// Payload refers to the JSON payload sent by a gateway or a server.
type Payload struct {
	Raw  []byte  `json:"-"`    // The raw unparsed response
	RXPK *[]RXPK `json:"rxpk"` // A list of RXPK messages transmitted if any
	Stat *Stat   `json:"stat"` // A Stat message transmitted if any
	TXPK *TXPK   `json:"txpk"` // A TXPK message transmitted if any
}

// Available packet commands
const (
	PUSH_DATA byte = iota // Sent by the gateway for an uplink message with data
	PUSH_ACK              // Sent by the gateway's recipient in response to a PUSH_DATA
	PULL_DATA             // Sent periodically by the gateway to keep a connection open
	PULL_RESP             // Sent by the gateway's recipient to transmit back data to the Gateway
	PULL_ACK              // Sent by the gateway's recipient in response to PULL_DATA
)

const VERSION = 0x01