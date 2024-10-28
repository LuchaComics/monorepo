import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";
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

import { GetImageFilePathFromDialog, GetVideoFilePathFromDialog, CreateNFT } from "../../../wailsjs/go/main/App";
import FormErrorBox from "../Reusable/FormErrorBox";
import FormInputField from "../Reusable/FormInputField";
import FormInputFieldWithButton from "../Reusable/FormInputFieldWithButton";
import FormTextareaField from "../Reusable/FormTextareaField";
import FormNFTMetadataAttributesField from "../Reusable/FormNFTMetadataAttributesField";


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

        console.log("onSubmitClick: Beginning...");
        setErrors({}); // Reset the errors in the GUI.

        // Variables used to create our new errors if we find them.
        let newErrors = {};
        let hasErrors = false;

        //////

        if (name === undefined || name === null || name === "") {
          newErrors["name"] = "missing value";
          hasErrors = true;
        }
        if (description === undefined || description === null || description === "") {
          newErrors["description"] = "missing value";
          hasErrors = true;
        }
        if (image === undefined || image === null || image === "") {
          newErrors["image"] = "missing value";
          hasErrors = true;
        }
        if (animation === undefined || animation === null || animation === "") {
          newErrors["animation"] = "missing value";
          hasErrors = true;
        }
        if (backgroundColor === undefined || backgroundColor === null || backgroundColor === "") {
          newErrors["backgroundColor"] = "missing value";
          hasErrors = true;
        }

        //////

        if (hasErrors) {
          console.log("onSubmitClick: Aboring because of error(s)");

          // Set the associate based error validation.
          setErrors(newErrors);

          // The following code will cause the screen to scroll to the top of
          // the page. Please see ``react-scroll`` for more information:
          // https://github.com/fisshy/react-scroll
          var scroll = Scroll.animateScroll;
          scroll.scrollToTop();

          return;
        }

        //////

        const attributesJSONString = JSON.stringify(attributes);

        // Submit the `dataDirectory` value to our backend.
        CreateNFT(name, description, image, animation, youtubeURL, externalURL, attributesJSONString, backgroundColor).then( (result) => {
            console.log("result:", result);
            // setForceURL("/startup")
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

                <FormNFTMetadataAttributesField
                  data={attributes}
                  onDataChange={setAttributes}
                />

                <FormInputField
                  label="Background Color"
                  name="backgroundColor"
                  placeholder=""
                  value={backgroundColor}
                  errorText={errors && errors.backgroundColor}
                  helpText=""
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
