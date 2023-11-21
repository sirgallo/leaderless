package connpool

import "errors"

import "google.golang.org/grpc"
import "google.golang.org/grpc/connectivity"
import "google.golang.org/grpc/credentials/insecure"
import "google.golang.org/grpc/encoding/gzip"


func NewConnectionPool(opts ConnectionPoolOpts) *ConnectionPool {
	return &ConnectionPool{
		maxConn: opts.MaxConn,
	}
}

func (cp *ConnectionPool) GetConnection(addr string, port string) (*grpc.ClientConn, error) {
	connections, loaded := cp.connections.Load(addr)
	if loaded {
		if len(connections.([]*grpc.ClientConn)) >= cp.maxConn { return nil, errors.New("max connections reached") }

		for _, conn := range connections.([]*grpc.ClientConn) {
			if conn != nil && conn.GetState() == connectivity.Ready { return conn, nil }
		}
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
	}

	newConn, connErr := grpc.Dial(addr + port, opts...)
	if connErr != nil { 
		cp.connections.Delete(addr)
		return nil, connErr 
	}

	emptyConns, loaded := cp.connections.LoadOrStore(addr, []*grpc.ClientConn{newConn})
	if loaded {
		connections := emptyConns.([]*grpc.ClientConn)
		cp.connections.Store(addr, append(connections, newConn))
	}
	
	return newConn, nil
}

func (cp *ConnectionPool) PutConnection(addr string, connection *grpc.ClientConn) (bool, error) {
	connections, loaded := cp.connections.Load(addr)
	if loaded {
		for _, conn := range connections.([]*grpc.ClientConn) {
			if conn == connection { return true, nil }
		}
	}

	closeErr := connection.Close()
	if closeErr != nil { return false, closeErr }
	
	return false, nil
}

func (cp *ConnectionPool) CloseConnections(addr string) (bool, error) {
	connections, loaded := cp.connections.Load(addr)
	if loaded {
		for _, conn := range connections.([]*grpc.ClientConn) {
			closeErr := conn.Close()
			if closeErr != nil { return false, closeErr }
		}
	}

	cp.connections.Delete(addr)
	return true, nil
}