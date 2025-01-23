from typing import Any
from boa3.builtin import public
from boa3.builtin.contract import Nep17TransferEvent, abort
from boa3.builtin.interop.contract import call_contract
from boa3.builtin.interop.runtime import calling_script_hash, check_witness, executing_script_hash
from boa3.builtin.interop.storage import get, put
from boa3.builtin.type import UInt160

# ------------------------------------------------------
#  Конфигурация токена
# ------------------------------------------------------
TOKEN_SYMBOL = 'MYNEO'
TOKEN_DECIMALS = 0

TOTAL_SUPPLY_KEY = b'totalSupply'
BALANCE_PREFIX = b'balance_'
ADMIN_KEY       = b'admin'   # владелец смарт-контракта

on_transfer = Nep17TransferEvent

@public
def symbol() -> str:
    return TOKEN_SYMBOL

@public
def decimals() -> int:
    return TOKEN_DECIMALS

@public
def totalSupply() -> int:
    return get(TOTAL_SUPPLY_KEY).to_int()

@public
def balanceOf(account: UInt160) -> int:
    return get(BALANCE_PREFIX + account).to_int()

@public
def transfer(from_addr: UInt160, to_addr: UInt160, amount: int, data: Any) -> bool:
    if amount < 0:
        abort()

    if not check_witness(from_addr):
        abort()

    if len(to_addr) != 20:
        abort()

    from_balance = balanceOf(from_addr)
    if from_balance < amount:
        abort()

    if from_addr != to_addr and amount != 0:
        put(BALANCE_PREFIX + from_addr, from_balance - amount)
        to_balance = balanceOf(to_addr)
        put(BALANCE_PREFIX + to_addr, to_balance + amount)

        on_transfer(from_addr, to_addr, amount)

        if data is not None:
            # Вызываем onNEP17Payment, если у целевого контракта есть такой метод
            call_contract(to_addr, 'onNEP17Payment', [from_addr, amount, data])

    return True

@public
def deploy(admin: UInt160):
    """
    При первом деплое: устанавливаем админа и минтим 100_000_000 токенов.
    """
    if totalSupply() != 0:
        return  # уже инициализировано

    put(ADMIN_KEY, admin)
    minted = 100_000_000
    put(TOTAL_SUPPLY_KEY, minted)
    put(BALANCE_PREFIX + admin, minted)
    on_transfer(None, admin, minted)
