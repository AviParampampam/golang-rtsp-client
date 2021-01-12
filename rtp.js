const net = require('net');

const server1 = net.createServer()
const server2 = net.createServer()

server1.on('connection', () => console.log('server1 connection'))
server2.on('connection', () => console.log('server2 connection'))

server1.listen(8000, '127.0.0.1')
server2.listen(8001, '127.0.0.1')

// ---

// const client = new net.Socket()
// client.connect(8000, '192.168.1.15', () => console.log("client connected1"))
// client.on('data', function(data) {
// 	console.log('Received: ' + data);
// 	client.destroy()
// });

// client.on('close', function() {
// 	console.log('Connection closed')
// });