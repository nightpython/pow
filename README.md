Algorithm Choice: Hashcash
The Word of Wisdom TCP Server incorporates the Hashcash algorithm to prevent abuse and mitigate Denial-of-Service (DoS) attacks. The Hashcash algorithm is a proof-of-work system designed to deter spam and other malicious activities by requiring a moderate amount of computational effort.

Here are the reasons why the Hashcash algorithm was chosen for this application:

Spam Prevention: By using Hashcash, the server ensures that clients requesting the "Word of Wisdom" service have performed a certain amount of computational work. This makes it costly and time-consuming for spammers to flood the server with excessive requests, as they would need to spend a significant amount of computational resources for each request.

DoS Mitigation: Hashcash helps to mitigate Denial-of-Service (DoS) attacks by making it computationally expensive for attackers to overwhelm the server with a large number of requests. By requiring a proof-of-work for each request, the server can prioritize legitimate clients who have invested computational effort over malicious actors.

Resource Fairness: The Hashcash algorithm ensures that clients need to contribute a reasonable amount of computational work to access the "Word of Wisdom" service. This helps distribute the server's resources more fairly among clients and prevents any single client from monopolizing the server's capacity.

Compatibility and Adoption: Hashcash is a well-known and widely adopted algorithm for proof-of-work systems. It has been used in various applications, including email spam filters and cryptocurrency mining. By leveraging an established algorithm, the Word of Wisdom TCP Server benefits from the existing research, security analysis, and community support surrounding Hashcash.

While Hashcash is an effective algorithm for preventing abuse and DoS attacks, it's important to note that the level of security it provides depends on the chosen parameters and the computational power available to the clients. It is advisable to periodically review and adjust the Hashcash parameters based on the specific needs and resources of your application.

Usage
To run the server and client, please refer to the instructions provided in the previous section.

Configuration
The server and client configurations are stored in the config/config.yaml file. You can modify this file to change the server and client settings as needed.