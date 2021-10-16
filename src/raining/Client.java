package raining;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.nio.channels.SocketChannel;

import com.fasterxml.jackson.databind.ObjectMapper;

public class Client {
	private ClientMsgProtocol protocol;
	private ObjectMapper mapper;
	private Gui gui;
	SocketChannel channel;
	
	public ObjectMapper getMapper() {
		return mapper;
	}

	public Gui getGui() {
		return gui;
	}

	public void setGui(Gui gui) {
		this.gui = gui;
	}

	public Client() {
		mapper = new ObjectMapper();
	}

	public void connect() {
		Thread t = new Thread(new Runnable() {

			@Override
			public void run() {
				try {
					SocketAddress addr = new InetSocketAddress("localhost", 5678);
					channel = SocketChannel.open(addr);
					protocol = new ClientMsgProtocol(Client.this, channel);
					protocol.readMsg();
				} catch (Exception e) {
					e.printStackTrace();
				} finally {
					if (channel != null) {
						try {
							channel.close();
						} catch (IOException e) {
							e.printStackTrace();
						}
					}
				}
			}

		});

		t.start();

	}

	public void sendMessage(JsonMsg msg) {
		try {
			protocol.writeMsg(mapper.writeValueAsBytes(msg));
		} catch (IOException e) {
			e.printStackTrace();
		}
	}

	public static void main(String[] args) throws IOException {
		Client client = new Client();
		client.connect();

		Gui gui = new Gui();

		client.setGui(gui);
		gui.setClient(client);
	}
}