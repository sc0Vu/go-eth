pragma solidity ^0.4.23;

import "./EIP20.sol";

contract HW is EIP20 {
	constructor() EIP20(10**(10+18), "Hello World", 18, "HW") public {
	}
}
