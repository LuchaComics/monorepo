export namespace big {
	
	export class Int {
	
	
	    static createFrom(source: any = {}) {
	        return new Int(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace domain {
	
	export class Validator {
	    id: string;
	    public_key_bytes: number[];
	
	    static createFrom(source: any = {}) {
	        return new Validator(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.public_key_bytes = source["public_key_bytes"];
	    }
	}
	export class BlockTransaction {
	    chain_id: number;
	    nonce_bytes: number[];
	    from?: number[];
	    to?: number[];
	    value: number;
	    tip: number;
	    data: number[];
	    type: string;
	    token_id_bytes: number[];
	    token_metadata_uri: string;
	    token_nonce_bytes: number[];
	    v_bytes: number[];
	    r_bytes: number[];
	    s_bytes: number[];
	    timestamp: number;
	    gas_price: number;
	    gas_units: number;
	
	    static createFrom(source: any = {}) {
	        return new BlockTransaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.chain_id = source["chain_id"];
	        this.nonce_bytes = source["nonce_bytes"];
	        this.from = source["from"];
	        this.to = source["to"];
	        this.value = source["value"];
	        this.tip = source["tip"];
	        this.data = source["data"];
	        this.type = source["type"];
	        this.token_id_bytes = source["token_id_bytes"];
	        this.token_metadata_uri = source["token_metadata_uri"];
	        this.token_nonce_bytes = source["token_nonce_bytes"];
	        this.v_bytes = source["v_bytes"];
	        this.r_bytes = source["r_bytes"];
	        this.s_bytes = source["s_bytes"];
	        this.timestamp = source["timestamp"];
	        this.gas_price = source["gas_price"];
	        this.gas_units = source["gas_units"];
	    }
	}
	export class BlockHeader {
	    chain_id: number;
	    number_bytes: number[];
	    prev_block_hash: string;
	    timestamp: number;
	    beneficiary: number[];
	    difficulty: number;
	    mining_reward: number;
	    state_root: string;
	    trans_root: string;
	    nonce_bytes: number[];
	    latest_token_id_bytes: number[];
	    tokens_root: string;
	
	    static createFrom(source: any = {}) {
	        return new BlockHeader(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.chain_id = source["chain_id"];
	        this.number_bytes = source["number_bytes"];
	        this.prev_block_hash = source["prev_block_hash"];
	        this.timestamp = source["timestamp"];
	        this.beneficiary = source["beneficiary"];
	        this.difficulty = source["difficulty"];
	        this.mining_reward = source["mining_reward"];
	        this.state_root = source["state_root"];
	        this.trans_root = source["trans_root"];
	        this.nonce_bytes = source["nonce_bytes"];
	        this.latest_token_id_bytes = source["latest_token_id_bytes"];
	        this.tokens_root = source["tokens_root"];
	    }
	}
	export class BlockData {
	    hash: string;
	    header?: BlockHeader;
	    header_signature_bytes: number[];
	    trans: BlockTransaction[];
	    validator?: Validator;
	
	    static createFrom(source: any = {}) {
	        return new BlockData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hash = source["hash"];
	        this.header = this.convertValues(source["header"], BlockHeader);
	        this.header_signature_bytes = source["header_signature_bytes"];
	        this.trans = this.convertValues(source["trans"], BlockTransaction);
	        this.validator = this.convertValues(source["validator"], Validator);
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
	
	
	export class NonFungibleTokenMetadataAttribute {
	    display_type: string;
	    trait_type: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new NonFungibleTokenMetadataAttribute(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.display_type = source["display_type"];
	        this.trait_type = source["trait_type"];
	        this.value = source["value"];
	    }
	}
	export class NonFungibleTokenMetadata {
	    image: string;
	    external_url: string;
	    description: string;
	    name: string;
	    attributes: NonFungibleTokenMetadataAttribute[];
	    background_color: string;
	    animation_url: string;
	    youtube_url: string;
	
	    static createFrom(source: any = {}) {
	        return new NonFungibleTokenMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.image = source["image"];
	        this.external_url = source["external_url"];
	        this.description = source["description"];
	        this.name = source["name"];
	        this.attributes = this.convertValues(source["attributes"], NonFungibleTokenMetadataAttribute);
	        this.background_color = source["background_color"];
	        this.animation_url = source["animation_url"];
	        this.youtube_url = source["youtube_url"];
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
	export class NonFungibleToken {
	    token_id?: big.Int;
	    metadata_uri: string;
	    metadata?: NonFungibleTokenMetadata;
	    state: string;
	
	    static createFrom(source: any = {}) {
	        return new NonFungibleToken(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token_id = this.convertValues(source["token_id"], big.Int);
	        this.metadata_uri = source["metadata_uri"];
	        this.metadata = this.convertValues(source["metadata"], NonFungibleTokenMetadata);
	        this.state = source["state"];
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
	
	
	export class Wallet {
	    label: string;
	    address?: number[];
	    keystore_bytes: number[];
	
	    static createFrom(source: any = {}) {
	        return new Wallet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.address = source["address"];
	        this.keystore_bytes = source["keystore_bytes"];
	    }
	}

}

