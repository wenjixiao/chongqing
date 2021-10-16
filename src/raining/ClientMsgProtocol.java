package raining;

import java.nio.channels.SocketChannel;

import javax.swing.SwingUtilities;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import raining.domain.Person;

public class ClientMsgProtocol extends AMsgProtocol {
	private Client client;

	public ClientMsgProtocol(Client client, SocketChannel channel) {
		super(channel);
		this.client = client;
	}

	public Client getClient() {
		return client;
	}

	@Override
	public void processMsg(byte[] data) throws Exception {
		ObjectMapper mapper = client.getMapper();
		JsonMsg msg = mapper.readValue(data, JsonMsg.class);
		
		switch (msg.getType()) {
		case 0:
			Person person = mapper.readValue(msg.getObject(), Person.class);
			System.out.println("got msg:"+person);
			break;
		default:
			System.out.println("hoho");
		}
	}
	
	private void println(String str) {
		SwingUtilities.invokeLater(new Runnable() {

			@Override
			public void run() {
				client.getGui().println(str);
			}
			
		});
	}

}
