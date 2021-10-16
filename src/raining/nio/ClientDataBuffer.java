package raining.nio;

import java.io.IOException;
import java.nio.channels.SocketChannel;

import raining.JsonMsg;
import raining.domain.Person;

public class ClientDataBuffer extends ADataReceiver {

	private SocketChannel channel;
	
	public ClientDataBuffer(SocketChannel channel) {
		super();
		this.channel = channel;
	}
	
	@Override
	public void processMsg(byte[] msgBody) {
		try {
			JsonMsg msg = Endecoder.getMapper().readValue(msgBody, JsonMsg.class);
			if(msg.getType()==0) {
				Person person = Endecoder.getMapper().readValue(msg.getObject(), Person.class);
				System.out.println("client get:"+person);
			}
		} catch (IOException e) {
			e.printStackTrace();
		}
	}

	public void writeMsg(byte[] msgBody) throws IOException {
		channel.write(pack2msg(msgBody));
	}
}
