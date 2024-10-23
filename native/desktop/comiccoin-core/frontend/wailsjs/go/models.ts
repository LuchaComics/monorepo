export namespace domain {
	
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

