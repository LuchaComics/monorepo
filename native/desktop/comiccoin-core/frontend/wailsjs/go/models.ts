export namespace domain {
	
	export class BlockTransaction {
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
	    // Go type: big
	    v?: any;
	    // Go type: big
	    r?: any;
	    // Go type: big
	    s?: any;
	    timestamp: number;
	    gas_price: number;
	    gas_units: number;
	
	    static createFrom(source: any = {}) {
	        return new BlockTransaction(source);
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
	        this.v = this.convertValues(source["v"], null);
	        this.r = this.convertValues(source["r"], null);
	        this.s = this.convertValues(source["s"], null);
	        this.timestamp = source["timestamp"];
	        this.gas_price = source["gas_price"];
	        this.gas_units = source["gas_units"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Token {
	    id: number;
	    owner?: number[];
	    metadata_uri: string;
	    nonce: number;
	
	    static createFrom(source: any = {}) {
	        return new Token(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.owner = source["owner"];
	        this.metadata_uri = source["metadata_uri"];
	        this.nonce = source["nonce"];
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

