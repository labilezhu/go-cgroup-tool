package gocgrouptool.stressjava;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import stress.Main;

@SpringBootApplication
public class StressJavaApplication {

	public static void main(String[] args) {
		Main.main(null);

		SpringApplication.run(StressJavaApplication.class, args);
	}

}
