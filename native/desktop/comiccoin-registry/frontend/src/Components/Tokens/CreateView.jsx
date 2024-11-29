import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import Scroll from "react-scroll";
import {
    faTasks,
    faGauge,
    faArrowRight,
    faUsers,
    faBarcode,
    faCubes,
    faCoins,
    faEllipsis,
    faPlus,
    faTimesCircle,
    faCheckCircle
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import { GetImageFilePathFromDialog, GetVideoFilePathFromDialog, CreateToken } from "../../../wailsjs/go/main/App";
import FormErrorBox from "../Reusable/FormErrorBox";
import FormInputField from "../Reusable/FormInputField";
import FormInputFieldWithButton from "../Reusable/FormInputFieldWithButton";
import FormTextareaField from "../Reusable/FormTextareaField";
import FormTokenMetadataAttributesField from "../Reusable/FormTokenMetadataAttributesField";


function CreateTokenView() {
    ////
    //// Component states.
    ////

    // --- GUI States ---

    const [forceURL, setForceURL] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const [errors, setErrors] = useState({});
    const [showCancelWarning, setShowCancelWarning] = useState(false);

    // --- Form States ---
    const [name, setName] = useState("");
    const [description, setDescription] = useState("");
    const [image, setImage] = useState("");
    const [animation, setAnimation] = useState("");
    const [youtubeURL, setYoutubeURL] = useState("");
    const [externalURL, setExternalURL] = useState("");
    const [attributes, setAttributes] = useState([]);
    const [backgroundColor, setBackgroundColor] = useState("");


    ////
    //// Event handling.
    ////
    const onSubmitClick = (e) => {
        e.preventDefault();

        // Reset the errors.
        setErrors({});

        // Update the GUI to let user know that the operation is under way.
        setIsLoading(true);

        const attributesJSONString = JSON.stringify(attributes);

        // Submit the `dataDirectory` value to our backend.
        CreateToken(name, description, image, animation, youtubeURL, externalURL, attributesJSONString, backgroundColor).then( (resp) => {
            console.log("result:", resp);
            setForceURL("/tokens")
        }).catch((errorJsonString)=>{
            console.log("errRes:", errorJsonString);

            let err = {};
            try {
                const errorObject = JSON.parse(errorJsonString);
                if (errorObject.name != "") {
                    err.name = errorObject.name;
                }
                if (errorObject.description != "") {
                    err.description = errorObject.description;
                }
                if (errorObject.image != "") {
                    err.image = errorObject.image;
                }
                if (errorObject.animation != "") {
                    err.animation = errorObject.animation;
                }
                if (errorObject.background_color != "") {
                    err.backgroundColor = errorObject.background_color;
                }
            } catch (e) {
                console.log("CreateToken:err:", e);
                err.message = errorJsonString;
            } finally {
                setErrors(err);

                // The following code will cause the screen to scroll to the top of
                // the page. Please see ``react-scroll`` for more information:
                // https://github.com/fisshy/react-scroll
                var scroll = Scroll.animateScroll;
                scroll.scrollToTop();
            }
        }).finally(() => {
            // this will be executed after then or catch has been executed
            console.log("promise has been resolved or rejected");

            // Update the GUI to let user know that the operation is completed.
            setIsLoading(false);
        });
    }

    ////
    //// API.
    ////



    ////
    //// Misc.
    ////

    useEffect(() => {
      let mounted = true;

      if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
      }

      return () => {
          mounted = false;
      };
    }, []);

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />;
    }

    if (isLoading) {
        return (
            <PageLoadingContent displayMessage="Submitting..." />
        );
    }

    return (
        <>
          <div class="container">
            <section class="section">
              <nav class="breadcrumb" aria-label="breadcrumbs">
                <ul>
                  <li>
                    <Link to="/dashboard" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Dashboard
                    </Link>
                  </li>
                  <li class="">
                    <Link to="/tokens" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faCubes} />
                      &nbsp;Tokens
                    </Link>
                  </li>
                  <li class="is-active">
                    <Link to="/tokens/new" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faPlus} />
                      &nbsp;New
                    </Link>
                  </li>
                </ul>
              </nav>

              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faPlus} />
                      &nbsp;New Token
                    </h1>
                  </div>
                </div>

                <FormErrorBox errors={errors} />

                <p class="pb-4">Please fill out all the required fields:</p>

                <FormInputField
                  label="Name"
                  name="name"
                  placeholder=""
                  value={name}
                  errorText={errors && errors.name}
                  helpText=""
                  onChange={(e) => setName(e.target.value)}
                  isRequired={true}
                  maxWidth="500px"
                />

                <FormTextareaField
                  label="Description"
                  name="description"
                  placeholder=""
                  value={description}
                  errorText={errors && errors.description}
                  helpText=""
                  onChange={(e) => setDescription(e.target.value)}
                  isRequired={true}
                  rows={6}
                />

                <FormInputFieldWithButton
                  label="Image"
                  name="image"
                  placeholder=""
                  value={image}
                  errorText={errors && errors.image}
                  helpText=""
                  onChange={(e) => setImage(e.target.value)}
                  isRequired={true}
                  maxWidth="500px"
                  buttonLabel={<><FontAwesomeIcon className="fas" icon={faEllipsis} /></>}
                  onButtonClick={(e) =>
                    GetImageFilePathFromDialog().then((imageRes) => {
                        if (imageRes !== "") {
                            setImage(imageRes);
                        }
                    })
                  }
                  inputOnlyDisabled={true}
                />

                <FormInputFieldWithButton
                  label="Animation"
                  name="animation"
                  placeholder=""
                  value={animation}
                  errorText={errors && errors.animation}
                  helpText=""
                  onChange={(e) => setAnimation(e.target.value)}
                  isRequired={true}
                  maxWidth="500px"
                  buttonLabel={<><FontAwesomeIcon className="fas" icon={faEllipsis} /></>}
                  onButtonClick={(e) =>
                    GetVideoFilePathFromDialog().then((animationRes) => {
                        if (animationRes !== "") {
                            setAnimation(animationRes);
                        }
                    })
                  }
                  inputOnlyDisabled={true}
                />

                <FormInputField
                  label="Youtube URL (Optional)"
                  name="youtubeURL"
                  placeholder=""
                  value={youtubeURL}
                  errorText={errors && errors.youtubeURL}
                  helpText=""
                  onChange={(e) => setYoutubeURL(e.target.value)}
                  isRequired={true}
                  maxWidth="500px"
                />

                <FormInputField
                  label="External URL (Optional)"
                  name="externalURL"
                  placeholder=""
                  value={externalURL}
                  errorText={errors && errors.externalURL}
                  helpText=""
                  onChange={(e) => setExternalURL(e.target.value)}
                  isRequired={true}
                  maxWidth="500px"
                />

                <FormTokenMetadataAttributesField
                  data={attributes}
                  onDataChange={setAttributes}
                />

                <FormInputField
                  label="Background Color"
                  name="backgroundColor"
                  placeholder="Ex: #FFFFFF"
                  value={backgroundColor}
                  errorText={errors && errors.backgroundColor}
                  helpText="Must be in a hexadecimal format"
                  onChange={(e) => setBackgroundColor(e.target.value)}
                  isRequired={true}
                  maxWidth="150px"
                />

              <div class="columns pt-5" style={{alignSelf: "flex-start"}}>
                <div class="column is-half ">
                  <button
                    class="button is-fullwidth-mobile"
                    onClick={(e) => setShowCancelWarning(true)}
                  >
                    <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                    &nbsp;Cancel
                  </button>
                </div>
                <div class="column is-half has-text-right">
                  <button
                    class="button is-primary is-fullwidth-mobile"
                    onClick={onSubmitClick}
                  >
                    <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                    &nbsp;Save
                  </button>
                </div>
              </div>


              </nav>
            </section>
          </div>
        </>
    )
}

export default CreateTokenView
