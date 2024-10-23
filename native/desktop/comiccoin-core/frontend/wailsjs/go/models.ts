export namespace domain {
	
	export class Transaction {
	    chain_id: number;
	    nonce: number;
	    from?: number[];
	    to?: number[];
	    value: number;
	    tip: number;
	    data: number[];
	    type: string;
	    token_id: number;
	    token_metadata_uri: string;
	    token_nonce: number;
	
	    static createFrom(source: any = {}) {
	        return new Transaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.chain_id = source["chain_id"];
	        this.nonce = source["nonce"];
	        this.from = source["from"];
	        this.to = source["to"];
	        this.value = source["value"];
	        this.tip = source["tip"];
	        this.data = source["data"];
	        this.type = source["type"];
	        this.token_id = source["token_id"];
	        this.token_metadata_uri = source["token_metadata_uri"];
	        this.token_nonce = source["token_nonce"];
	    }
	}
	export class Wallet {
	    label: string;
	    address?: number[];
	    filepath: string;
	
	    static createFrom(source: any = {}) {
	        return new Wallet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.address = source["address"];
	        this.filepath = source["filepath"];
	    }
	}

}

