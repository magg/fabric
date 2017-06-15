package interceptor



import (
	"github.com/hyperledger/fabric/client"
)


  //const REQUEST_KEY = "x-request-id"
  //const TRACE_KEY = "x-b3-traceid"

  var headers = []string{
      "x-ot-span-context",
      "x-request-id",
      "x-b3-traceid",
      "x-b3-spanid",
      "x-b3-parentspanid",
      "x-b3-sampled",
      "x-b3-flags",
    }



    var hm = map[string]string{}

  // Returns a slice of strings suitable for passing to grpc/metadata.Pairs
  func GRPCMetadata() []string {
  	r := client.GetRPCMetadata()
  	return []string{headers[1],r.RequestID,
      headers[2], r.TraceID,
      headers[3], r.SpanID,
      headers[4], r.ParentSpanID,
      headers[5], r.SampledID,
      headers[6], r.FlagsID,
      headers[0], r.OtSpan,

    }
  }


func GRPCRecieved(md map[string][]string) {


getIDs(md)

/*

	req_strs, ok := md[REQUEST_KEY]
	if !ok || len(req_strs) < 1 {
		return
	}
	req := getID(md, REQUEST_KEY)

	trace_strs, ok := md[TRACE_KEY]
	if !ok || len(trace_strs) < 1 {
		return
	}
	traceID := getID(md, TRACE_KEY)

*/

getIDs(md)

	client.RPCReceived(client.RPCMetadata{
		RequestID: hm[headers[1]],
		TraceID: hm[headers[2]],
    SpanID:  hm[headers[3]],
    ParentSpanID:  hm[headers[4]],
    SampledID:   hm[headers[5]],
    FlagsID:  hm[headers[6]],
    OtSpan:  hm[headers[0]],
	})
}

func GRPCReturned(md map[string][]string) {

  /*
  req_strs, ok := md[REQUEST_KEY]
  if !ok || len(req_strs) < 1 {
    return
  }
  req := getID(md, REQUEST_KEY)

  trace_strs, ok := md[TRACE_KEY]
  if !ok || len(trace_strs) < 1 {
    return
  }
  traceID := getID(md, TRACE_KEY)
*/
	client.RPCReturned(client.RPCMetadata{
    RequestID: hm[headers[1]],
		TraceID: hm[headers[2]],
    SpanID:  hm[headers[3]],
    ParentSpanID:  hm[headers[4]],
    SampledID:   hm[headers[5]],
    FlagsID:  hm[headers[6]],
    OtSpan:  hm[headers[0]],

	})
}

func getIDs(md map[string][]string) {


		for i, n := range md {
		        fmt.Printf("OLD %s: %s\n", i, n)
		    }


    for i := 0; i < len(headers); i++ {
      if id := getID(md, headers[i]); len(id) > 0 {
				fmt.Printf("key %s bucket %s\n", headers[i], id)
        hm[headers[i]] = id
      }
    }

}


// getID parses an id from the metadata.
func getID(md map[string][]string, name string) string {
	for _, str := range md[name] {

		if len(str) > 0 {
			return str
		}
	}
	return ""
}
