package raining;

import java.awt.BorderLayout;
import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;
import java.awt.event.WindowAdapter;
import java.awt.event.WindowEvent;
import java.io.IOException;
import java.util.Arrays;
import java.util.List;

import javax.swing.JFrame;
import javax.swing.JTextArea;
import javax.swing.JTextField;
import javax.swing.event.AncestorEvent;
import javax.swing.event.AncestorListener;

import com.fasterxml.jackson.core.JsonProcessingException;

import raining.domain.Person;
import raining.messages.Login;

public class Gui {
	private Client client;
	private JTextArea output;
	private JTextField input;

	public Client getClient() {
		return client;
	}

	public void setClient(Client client) {
		this.client = client;
	}

	class MyWindowListener extends WindowAdapter {
		@Override
		public void windowClosing(WindowEvent e) {
			System.out.println("ssslslsllssll");
		}
	}

	public Gui() throws IOException {
		JFrame frame = new JFrame("my client");

		frame.addWindowListener(new MyWindowListener());
		frame.setSize(400, 200);
		frame.setLayout(new BorderLayout());

		output = new JTextArea();
		frame.add(output, BorderLayout.CENTER);

		input = new JTextField();

		input.addAncestorListener(new AncestorListener() {

			@Override
			public void ancestorAdded(AncestorEvent event) {
				input.requestFocus();
			}

			@Override
			public void ancestorRemoved(AncestorEvent event) {
			}

			@Override
			public void ancestorMoved(AncestorEvent event) {
			}

		});

		input.addActionListener(new ActionListener() {

			@Override
			public void actionPerformed(ActionEvent e) {
				if (e.getSource() == input) {
					String cmd = input.getText().trim();
					try {
						parseCmd(cmd);
					} catch (JsonProcessingException e1) {
						// TODO Auto-generated catch block
						e1.printStackTrace();
					}
					input.setText("");
				}
			}

		});

		frame.add(input, BorderLayout.SOUTH);

		frame.setLocationRelativeTo(null);
		frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
		frame.setVisible(true);
	}

	public void println(String s) {
		output.append(s);
		output.append("\n");
	}

	public void parseCmd(String cmd) throws JsonProcessingException {
		List<String> words = Arrays.asList(cmd.split(" "));
		System.out.println("cmd string:" + words);
		if (!words.isEmpty()) {
			String command = words.get(0);
			switch (command) {
			case "login":
				Login login = new Login();
				login.setPid(words.get(1));
				login.setPassword(words.get(2));

				
//				Message<Login> msg = new Message<Login>();
//				msg.setType(Message.Type.Login);
//				msg.setMsg(login);

//				client.sendMessage(msg);

				break;
			case "person":
				Person person = new Person();
				person.setName("wenjixiao");
				person.setYearold((byte)20);
				
				JsonMsg msg = new JsonMsg();
				msg.setType((byte)0);
				msg.setObject(client.getMapper().writeValueAsBytes(person));
				client.sendMessage(msg);
				
			default:
				System.out.println("***unknow command***");
			}

		}

	}

}
