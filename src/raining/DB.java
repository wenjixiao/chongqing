package raining;

import raining.domain.Player;

public class DB {
	public static Player getPlayer(String pid,String password) {
		Player p = new Player();
		p.setPid("wenjixiao");
		p.setLevel("3k");
		p.setPassword("123");
		return p;
	}
}
