// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {domain} from '../models';

export function CreateWallet(arg1:string,arg2:string,arg3:string):Promise<string>;

export function DefaultWalletAddress():Promise<string>;

export function GetDataDirectoryFromDialog():Promise<string>;

export function GetDataDirectoryFromPreferences():Promise<string>;

export function GetDefaultDataDirectory():Promise<string>;

export function GetIsBlockhainNodeRunning():Promise<boolean>;

export function GetRecentTransactions(arg1:string):Promise<Array<domain.BlockTransaction>>;

export function GetTotalCoins(arg1:string):Promise<number>;

export function GetTotalTokens(arg1:string):Promise<number>;

export function ListWallets():Promise<Array<domain.Wallet>>;

export function SaveDataDirectory(arg1:string):Promise<void>;

export function ShutdownApp():Promise<void>;
