package raining.nio;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.SocketChannel;

import raining.JsonMsg;
import raining.domain.Person;

public class NioClient {
	public static void main(String[] args) {
		NioClient client = new NioClient();
		client.connect("localhost", 20000);
	}

	public void connect(String serverIp, int serverPort) {
		try {
			SocketChannel channel = SocketChannel.open(new InetSocketAddress(serverIp, serverPort));
			
			Person person = new Person();
			person.setName("wenjixiao");
			person.setYearold((byte)20);
			
			JsonMsg msg = new JsonMsg();
			msg.setType((byte)0);
			msg.setObject(Endecoder.getMapper().writeValueAsBytes(person));
			
			byte[] mydata = Endecoder.getMapper().writeValueAsBytes(msg);
			ClientDataBuffer dataBuffer = new ClientDataBuffer(channel);
			dataBuffer.writeMsg(mydata);
			
			//read now!
			ByteBuffer buffer = ByteBuffer.allocate(1024);
			channel.read(buffer);
			buffer.flip();
			
			dataBuffer.receiveData(buffer);
			Thread.sleep(6000);
			channel.close();
			System.out.println("----client exit----");
		} catch (IOException e) {
			e.printStackTrace();
		} catch (InterruptedException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}

	}
}
