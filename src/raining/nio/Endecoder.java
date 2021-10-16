package raining.nio;

import com.fasterxml.jackson.databind.ObjectMapper;

public class Endecoder {
	private static ObjectMapper mapper;
	
	static {
		mapper = new ObjectMapper();
	}
	
	public static ObjectMapper getMapper() {
		return mapper;
	}
}
