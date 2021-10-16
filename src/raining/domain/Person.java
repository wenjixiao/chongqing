package raining.domain;

public class Person {
	private String name;
	private byte yearold;
	
	public String getName() {
		return name;
	}
	public void setName(String name) {
		this.name = name;
	}
	public byte getYearold() {
		return yearold;
	}
	public void setYearold(byte yearold) {
		this.yearold = yearold;
	}
	@Override
	public String toString() {
		return "Person [name=" + name + ", yearold=" + yearold + "]";
	}
	
	
}
