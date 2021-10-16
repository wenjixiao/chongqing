package raining.nio;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.nio.channels.ServerSocketChannel;
import java.nio.channels.SocketChannel;
import java.util.Iterator;

public class NioServer {
	private static final int READ_BUF_SIZE = 1024;

	public NioServer() {
	}

	public void startServer(String serverIp, int serverPort) throws IOException {
		ServerSocketChannel serviceChannel = ServerSocketChannel.open();
		InetSocketAddress localAddr = new InetSocketAddress(serverIp, serverPort);
		serviceChannel.bind(localAddr);
		serviceChannel.configureBlocking(false);

		Selector selector = Selector.open();
		serviceChannel.register(selector, SelectionKey.OP_ACCEPT);

		while (true) {
			selector.select();
			Iterator<SelectionKey> keyIterator = selector.selectedKeys().iterator();
			while (keyIterator.hasNext()) {
				SelectionKey key = keyIterator.next();
				keyIterator.remove();
				try {
					if (key.isAcceptable()) {
						System.out.println("accepted...");
						ServerSocketChannel server = (ServerSocketChannel) key.channel();
						SocketChannel channel = server.accept();
						channel.configureBlocking(false);
						ServerDataBuffer dataBuffer = new ServerDataBuffer();
						channel.register(selector, SelectionKey.OP_READ, dataBuffer);
					} else if (key.isReadable()) {
//						System.out.println("readable...");
						SocketChannel channel = (SocketChannel) key.channel();
						ServerDataBuffer dataBuffer = (ServerDataBuffer) key.attachment();
						
						ByteBuffer buffer = ByteBuffer.allocate(READ_BUF_SIZE);
						channel.read(buffer);
						buffer.flip();
						dataBuffer.receiveData(buffer);

						channel.register(selector, SelectionKey.OP_WRITE, dataBuffer);
					} else if (key.isWritable()) {
						// Basically,always can write? So,what we care is sendBuf's length.
//						System.out.println("writable...");
						SocketChannel channel = (SocketChannel) key.channel();
						ServerDataBuffer dataBuffer = (ServerDataBuffer) key.attachment();
						
						ByteBuffer sendBuf = dataBuffer.getSendBuf();
						if (sendBuf.hasRemaining()) {
							channel.write(sendBuf.flip());
							sendBuf.clear();//ready for write again
						}
						
						channel.register(selector, SelectionKey.OP_READ, dataBuffer);
					}
				} catch (IOException e) {
					key.cancel();
					key.channel().close();

					e.printStackTrace();
					System.out.println("****exit****");
				}

			}
		}
	}
	
	public static void main(String[] args) {
		NioServer myserver = new NioServer();
		try {
			myserver.startServer("localhost", 20000);
		} catch (IOException e) {
			e.printStackTrace();
		}
	}

}
