package client

import (

	"github.com/hyperledger/fabric/local"
)

var token local.Token

type localStorage struct {
	requestID        string
	traceID          string
  spanID           string
  parentSpanID     string
  sampledID        string
  flagsID          string
  otSpan           string
  //redundancies []string
	//tags         []string
}

// exported type for RPC calls
type RPCMetadata struct {
  RequestID        string
	TraceID          string
  SpanID           string
  ParentSpanID     string
  SampledID        string
  FlagsID          string
  OtSpan           string
}

func init() {
	token = local.Register(&localStorage{
    requestID:        "",
  	traceID:          "",
    spanID:           "",
    parentSpanID:     "",
    sampledID:        "",
    flagsID:          "",
    otSpan:           "",

  }, local.Callbacks{
		func(l interface{}) interface{} {
			// deep copy l
			n := *(l.(*localStorage))
			return &n
		},
	})
}

// Runs the given function in a new goroutine, but copies the
// local vars from the current goroutine first.
func XGo(f func()) {
	go func(f1 func(), f2 func()) {
		f1()
		f2()
	}(local.GetSpawnCallback(), f)
}

func getLocal() *localStorage {
	return local.GetLocal(token).(*localStorage)
}

func GetRPCMetadata() RPCMetadata {
	l := getLocal()
	var r RPCMetadata
  r.RequestID    = l.requestID
  r.TraceID      = l.traceID
  r.SpanID       = l.spanID
  r.ParentSpanID = l.parentSpanID
  r.SampledID    = l.sampledID
  r.FlagsID      = l.flagsID
  r.OtSpan       = l.otSpan
	return r
}

func (r *RPCMetadata) Set() {
	l := getLocal()
  r.RequestID    = l.requestID
  r.TraceID      = l.traceID
  r.SpanID       = l.spanID
  r.ParentSpanID = l.parentSpanID
  r.SampledID    = l.sampledID
  r.FlagsID      = l.flagsID
  r.OtSpan       = l.otSpan

}

func RPCReceived(r RPCMetadata) {
  SetTraceID(r.TraceID)
	SetRequestID(r.RequestID)
  SetSpanID(r.SpanID)
  SetParentSpanID(r.ParentSpanID)
  SetSampledID(r.SampledID)
  SetFlagsID(r.FlagsID)
  SetOtSpan(r.OtSpan)

	//(msg)
}

func RPCReturned(r RPCMetadata) {
	SetTraceID(r.TraceID)
	SetRequestID(r.RequestID)
  SetSpanID(r.SpanID)
  SetParentSpanID(r.ParentSpanID)
  SetSampledID(r.SampledID)
	//Log(msg)
}



func SetSpanID(spanID string) {
	getLocal().spanID = spanID
}

func SetParentSpanID(parentSpanID string) {
	getLocal().parentSpanID = parentSpanID
}
func SetSampledID(sampledID string) {
	getLocal().sampledID = sampledID
}

func SetFlagsID(flagsID string) {
	getLocal().flagsID = flagsID
}
func SetOtSpan(otSpan string) {
	getLocal().otSpan = otSpan
}



// SetEventID sets the current goroutine's X-Trace Event ID.
// This should be used when propagating Event IDs over RPC
// calls or other channels.
//
// WARNING: This will overwrite any previous Event ID,
// so call with caution.
func SetTraceID(traceID string) {
	getLocal().traceID = traceID
}

// SetTaskID sets the current goroutine's X-Trace Task ID.
// This should be used when propagating Task IDs over RPC
// calls or other channels.
//
// WARNING: This will overwrite any previous Task ID,
// so call with caution.
func SetRequestID(requestID string) {
	getLocal().requestID = requestID
}


/*
func AddTags(str ...string) {
	if getLocal().tags == nil {
		getLocal().tags = str
	} else {
		getLocal().tags = append(getLocal().tags, str...)
	}
}

func NewTask(tags ...string) {
	SetTaskID(randInt64())
	SetEventID(randInt64())
	getLocal().tags = tags
}
*/
// GetEventID gets the current goroutine's X-Trace Event ID.
// Note that if one has not been set yet, GetEventID will
// return 0. This should be used when propagating Event IDs
// over RPC calls or other channels.
func GetTraceID() (traceID string) {
	return getLocal().traceID
}

// GetTaskID gets the current goroutine's X-Trace Task ID.
// Note that if one has not been set yet, GetTaskID will
// return 0. This should be used when propagating Task IDs
// over RPC calls or other channels.
func GetRequestID() (requestID string) {
	return getLocal().requestID
}
/*
func AddRedundancies(eventIDs ...int64) {
	getLocal().redundancies = append(getLocal().redundancies, eventIDs...)
}

func PopRedundancies() []int64 {
	eventIDs := append([]int64{}, getLocal().redundancies...)
	getLocal().redundancies = []int64{}
	return eventIDs
}

func newEvent() (parent, event int64) {
	parent = GetEventID()
	event = randInt64()
	SetEventID(event)
	return parent, event
}

func randInt64() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(fmt.Errorf("could not read random bytes: %v", err))
	}
	// shift to guarantee high bit is 0 and thus
	// int64 version is non-negative
	return int64(binary.BigEndian.Uint64(b[:]) >> 1)
}
*/
