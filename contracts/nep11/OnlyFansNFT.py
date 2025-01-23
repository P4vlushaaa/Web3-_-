from typing import Any, Dict, List
from boa3.builtin import public
from boa3.builtin.contract import abort
from boa3.builtin.interop.runtime import check_witness
from boa3.builtin.interop.storage import get, put, delete, find, FindOptions
from boa3.builtin.type import UInt160

TOKEN_PREFIX   = b'token_'
OWNER_PREFIX   = b'owner_'
SUPPLY_KEY     = b'totalSupply'
ADMIN_KEY      = b'admin'

@public
def symbol() -> str:
    return "OFNFT"  # OnlyFans NFT

@public
def name() -> str:
    return "Web3OnlyFansNFT"

@public
def totalSupply() -> int:
    return get(SUPPLY_KEY).to_int()

@public
def balanceOf(owner: UInt160) -> int:
    return len(_getTokensOf(owner))

@public
def ownerOf(token_id: str) -> UInt160:
    token_data = get(TOKEN_PREFIX + token_id.encode())
    if not token_data:
        abort()
    info = _deserialize_token(token_data)
    return info['owner']

@public
def tokensOf(owner: UInt160) -> List[str]:
    return _getTokensOf(owner)

@public
def properties(token_id: str) -> Dict[str, Any]:
    token_data = get(TOKEN_PREFIX + token_id.encode())
    if not token_data:
        abort()
    return _deserialize_token(token_data)

@public
def transfer(to: UInt160, token_id: str, data: Any) -> bool:
    if len(to) != 20:
        abort()

    token_info = _readToken(token_id)
    from_addr = token_info['owner']

    if not check_witness(from_addr):
        abort()

    if from_addr == to:
        return True

    _removeTokenOf(from_addr, token_id)
    _addTokenOf(to, token_id)

    token_info['owner'] = to
    put(TOKEN_PREFIX + token_id.encode(), _serialize_token(token_info))
    return True

@public
def mint(token_id: str, name: str, blurred_ref: str, full_ref: str, price: int) -> bool:
    admin = get(ADMIN_KEY)
    if not admin or len(admin) != 20:
        abort()

    if not check_witness(admin.to_uint160()):
        abort()

    existing = get(TOKEN_PREFIX + token_id.encode())
    if existing:
        abort()  # уже есть

    owner = admin.to_uint160()

    info = {
        "owner": owner,
        "name": name,
        "blurred_ref": blurred_ref,
        "full_ref": full_ref,
        "price": price
    }
    put(TOKEN_PREFIX + token_id.encode(), _serialize_token(info))

    current_supply = totalSupply()
    put(SUPPLY_KEY, current_supply + 1)
    _addTokenOf(owner, token_id)

    return True

@public
def deploy(admin: UInt160):
    if totalSupply() != 0:
        return
    put(ADMIN_KEY, admin)

def _readToken(token_id: str) -> Dict[str, Any]:
    token_data = get(TOKEN_PREFIX + token_id.encode())
    if not token_data:
        abort()
    return _deserialize_token(token_data)

def _getTokensOf(owner: UInt160) -> List[str]:
    prefix = OWNER_PREFIX + owner
    result: List[str] = []
    iterator = find(prefix, options=FindOptions.KeysOnly)
    while iterator.next():
        key = iterator.value
        token_id = key[len(prefix):].decode()
        result.append(token_id)
    return result

def _addTokenOf(owner: UInt160, token_id: str):
    owner_key = OWNER_PREFIX + owner + token_id.encode()
    put(owner_key, b'1')

def _removeTokenOf(owner: UInt160, token_id: str):
    owner_key = OWNER_PREFIX + owner + token_id.encode()
    delete(owner_key)

def _serialize_token(token_info: Dict[str, Any]) -> bytes:
    from boa3.builtin import serialize
    return serialize(token_info)

def _deserialize_token(data: bytes) -> Dict[str, Any]:
    from boa3.builtin import deserialize
    return deserialize(data)
