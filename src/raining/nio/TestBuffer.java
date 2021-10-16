package raining.nio;

import java.nio.IntBuffer;

public class TestBuffer {

	public static void main(String[] args) {
		IntBuffer nums = IntBuffer.allocate(10);
		for(int i=1;i<=6;i++) {
			nums.put(i);
		}
		print(nums);
		nums.flip();
		print(nums);
		print(nums.get());
		print(nums.get());
		print(nums);
		print(nums.compact());
	}

	public static void print(IntBuffer b) {
		System.out.println(b);
	}
	
	public static void print(int i) {
		System.out.println(i);
	}
	
}
