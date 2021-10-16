package raining.nio;

import java.io.IOException;
import java.nio.ByteBuffer;

import com.fasterxml.jackson.core.JsonProcessingException;

import raining.JsonMsg;
import raining.domain.Person;

public class ServerDataBuffer extends ADataReceiver {
	private static final int SEND_BUF_LEN = 1024 * 3;

	private ByteBuffer sendBuf;
	
	public ServerDataBuffer() {
		super();
		sendBuf = ByteBuffer.allocate(SEND_BUF_LEN);
	}
	
	@Override
	public void processMsg(byte[] data) {
		System.out.println("hehll");
		try {
			JsonMsg msg = Endecoder.getMapper().readValue(data, JsonMsg.class);
			switch (msg.getType()) {
			case 0:
				Person person = Endecoder.getMapper().readValue(msg.getObject(), Person.class);
				System.out.println("got msg:"+person);
				person.setYearold((byte) 40);
				JsonMsg jsonMsg = new JsonMsg();
				jsonMsg.setType((byte) 0);
				jsonMsg.setObject(Endecoder.getMapper().writeValueAsBytes(person));
				writeMsg(jsonMsg);
				break;
			default:
				System.out.println("hoho");
			}
		} catch (IOException e) {
			e.printStackTrace();
		}
	}

	public void writeMsg(JsonMsg msg) throws JsonProcessingException {
		sendBuf.put(pack2msg(Endecoder.getMapper().writeValueAsBytes(msg)));
	}
	
	public ByteBuffer getSendBuf() {
		return sendBuf;
	}
	
}
