package raining;

public class JsonMsg {
	private byte type;
	private byte[] object;
	
	public byte getType() {
		return type;
	}
	public void setType(byte head) {
		this.type = head;
	}
	public byte[] getObject() {
		return object;
	}
	public void setObject(byte[] body) {
		this.object = body;
	}
	
}
