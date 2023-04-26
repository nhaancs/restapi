CREATE TABLE `users` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(255) NOT NULL UNIQUE,
  `hashed_password` varchar(255) NOT NULL,
  `salt` varchar(255) NOT NULL
);


