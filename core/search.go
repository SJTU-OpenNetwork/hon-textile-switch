package core

import (
	"sync"

//	"github.com/segmentio/ksuid"
//	"github.com/SJTU-OpenNetwork/hon-textile-switch/broadcast"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/pb"
//	"github.com/SJTU-OpenNetwork/hon-textile-switch/util"
)

// queryResultSet holds a unique set of search results
type queryResultSet struct {
	options *pb.QueryOptions
	items   map[string]*pb.QueryResult
	lock    sync.Mutex
}

// newQueryResultSet returns a new queryResultSet
func newQueryResultSet(options *pb.QueryOptions) *queryResultSet {
	return &queryResultSet{
		options: options,
		items:   make(map[string]*pb.QueryResult, 0),
	}
}

// Add only adds a result to the set if it's newer than last
func (s *queryResultSet) Add(items ...*pb.QueryResult) []*pb.QueryResult {
	s.lock.Lock()
	defer s.lock.Unlock()

	var added []*pb.QueryResult
	for _, i := range items {
		//last := s.items[i.Id]
		s.items[i.Id] = i
		added = append(added, i)
	}

	return added
}

// List returns the items as a slice
func (s *queryResultSet) List() []*pb.QueryResult {
	s.lock.Lock()
	defer s.lock.Unlock()

	var list []*pb.QueryResult
	for _, i := range s.items {
		list = append(list, i)
	}

	return list
}

// Full returns whether or not the number of results meets or exceeds limit
func (s *queryResultSet) Full() bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.items) >= int(s.options.Limit)
}


// Search searches the network based on the given query
//func (t *Textile) searchByPubsub(query *pb.Query) (<-chan *pb.QueryResult, <-chan error, *broadcast.Broadcaster) {
//	return nil,nil,nil
//    /*
//    query = queryDefaults(query)
//	query.Id = ksuid.New().String()
//
//	var searchChs []chan *pb.QueryResult
//
//    // remote results channel(s)
//	clientCh := make(chan *pb.QueryResult)
//	searchChs = append(searchChs, clientCh)
//
//	resultCh := mergeQueryResults(searchChs)
//	errCh := make(chan error)
//	cancel := broadcast.NewBroadcaster(0)
//
//	go func() {
//		defer func() {
//			for _, ch := range searchChs {
//				close(ch)
//			}
//		}()
//		results := newQueryResultSet(query.Options)
//
//		canceler := cancel.Listen()
//		if err := t.cafe.searchPubSub(query, func(res *pb.QueryResults) bool {
//			for _, n := range results.Add(res.Items...) {
//				clientCh <- n
//			}
//			return results.Full()
//		}, canceler.Ch, false); err != nil {
//			errCh <- err
//			return
//		}
//	}()
//
//	return resultCh, errCh, cancel
//    */
//}
//
//// mergeQueryResults merges results from mulitple queries
//func mergeQueryResults(cs []chan *pb.QueryResult) chan *pb.QueryResult {
//	out := make(chan *pb.QueryResult)
//	var wg sync.WaitGroup
//	wg.Add(len(cs))
//	for _, c := range cs {
//		go func(c chan *pb.QueryResult) {
//			for v := range c {
//				out <- v
//			}
//			wg.Done()
//		}(c)
//	}
//	go func() {
//		wg.Wait()
//		close(out)
//	}()
//	return out
//}
