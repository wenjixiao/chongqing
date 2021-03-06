package raining.domain;

import java.io.Serializable;

public class Stone implements Serializable {
	private static final long serialVersionUID = 1L;

	public byte getX() {
		return x;
	}

	public void setX(byte x) {
		this.x = x;
	}

	public byte getY() {
		return y;
	}

	public void setY(byte y) {
		this.y = y;
	}

	public Color getColor() {
		return color;
	}

	public void setColor(Color color) {
		this.color = color;
	}

	public boolean isPass() {
		return pass;
	}

	public void setPass(boolean pass) {
		this.pass = pass;
	}

	private byte x;
	private byte y;
	private Color color;
	private boolean pass;

	@Override
	public String toString() {
		return "Stone [x=" + x + ", y=" + y + ", color=" + color + ", pass=" + pass + "]";
	}

}
