# Markov Passwords

As discussed in my term paper, this program reads [plain text password leaks](https://hashes.org/leaks.php) and creates a transistion matrix. How many previous states are recorded, i.e. the order of the Markov chain, is given by the user with the `-n` flag. If you need further options, look at `markov analyze -h`

The result may then be used to generate new passwords. To see all options, look at `markov generate -h`.
