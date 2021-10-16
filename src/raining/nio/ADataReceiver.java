package raining.nio;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

public abstract class ADataReceiver {
	public static final int MSG_HEADER_LEN = 4;
	private static final int RECEIVE_BUF_LEN = 1024 * 6;
	
	private int bodyLen;
	private ByteBuffer receiveBuf;
	
	public ADataReceiver() {
		bodyLen = 0;
		receiveBuf = ByteBuffer.allocate(RECEIVE_BUF_LEN);
	}
	
	public void receiveData(ByteBuffer buffer) {
		if(buffer.hasRemaining()) {
			System.out.println("received data:"+buffer);
		}
		receiveBuf.put(buffer);
		receiveBuf.flip();

		while (receiveBuf.hasRemaining()) {
			int n = receiveBuf.remaining();
			if (bodyLen == 0) {
				if (n >= MSG_HEADER_LEN) {
					byte[] header = new byte[MSG_HEADER_LEN];
					receiveBuf = receiveBuf.get(header);
					ByteBuffer headBuf = ByteBuffer.wrap(header).order(ByteOrder.LITTLE_ENDIAN);
					bodyLen = headBuf.getInt();
				} else {
					break;
				}
			}

			if (bodyLen > 0) {
				if (receiveBuf.remaining() >= bodyLen) {
					byte[] body = new byte[bodyLen];
					receiveBuf = receiveBuf.get(body);
					processMsg(body);
					bodyLen = 0;
				}
				break;
			}
		}

		receiveBuf.compact();
	}
	
	public static ByteBuffer pack2msg(byte[] msgBody) {
		int bodyLen = msgBody.length;
		ByteBuffer msgBuf = ByteBuffer.allocate(MSG_HEADER_LEN + bodyLen);

		ByteBuffer headBuf = ByteBuffer.allocate(MSG_HEADER_LEN).order(ByteOrder.LITTLE_ENDIAN);
		headBuf.putInt(bodyLen);
		headBuf.flip();
		
		msgBuf.put(headBuf);
		msgBuf.put(msgBody);
		msgBuf.flip();
		
		return msgBuf;
	}

	public abstract void processMsg(byte[] msgBody);

}
