export namespace domain {
	
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

}

