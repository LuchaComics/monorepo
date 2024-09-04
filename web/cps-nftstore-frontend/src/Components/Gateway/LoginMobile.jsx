import React from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowRight,
  faEnvelope,
  faKey,
  faTriangleExclamation,
} from "@fortawesome/free-solid-svg-icons";

import FormErrorBox from "../Reusable/FormErrorBox";

const LoginMobile = ({
  errors,
  isUnauthorized,
  validation,
  email,
  onEmailChange,
  password,
  onPasswordChange,
  onButtonClick,
  rememberMe,
  setRememberMe,
}) => {
  return (
    <div class="box is-rounded">
      {/* Start Logo */}
      <nav class="level">
        <div class="level-item has-text-centered">
          <figure class="image">
            <img
              src="/img/cps-logo.jpeg"
              style={{ width: "110px" }}
              alt="CPS Logo"
            />
          </figure>
        </div>
      </nav>
      {/* End Logo */}
      <form>
        <h1 className="title is-4 has-text-centered">Sign In</h1>
        {isUnauthorized === "true" && (
          <article class="message is-danger">
            <div class="message-body">
              <FontAwesomeIcon className="fas" icon={faTriangleExclamation} />
              &nbsp;Your session has ended.
              <br />
              Please login again
            </div>
          </article>
        )}
        <FormErrorBox errors={errors} />

        <div class="field">
          <label class="label is-small has-text-grey-darker">Email</label>
          <div class="control has-icons-left has-icons-right">
            <input
              class={`input ${errors && errors.email && "is-danger"} ${
                validation && validation.email && "is-success"
              }`}
              type="email"
              placeholder="Email"
              value={email}
              onChange={onEmailChange}
            />
            <span class="icon is-small is-left">
              <FontAwesomeIcon className="fas" icon={faEnvelope} />
            </span>
          </div>
          {errors && errors.email && (
            <p class="help is-danger">{errors.email}</p>
          )}
        </div>

        <div class="field">
          <label class="label is-small has-text-grey-darker">Password</label>
          <div class="control has-icons-left has-icons-right">
            <input
              class={`input ${errors && errors.password && "is-danger"} ${
                validation && validation.password && "is-success"
              }`}
              type="password"
              placeholder="Password"
              value={password}
              onChange={onPasswordChange}
            />
            <span class="icon is-small is-left">
              <FontAwesomeIcon className="fas" icon={faKey} />
            </span>
          </div>
          {errors && errors.password && (
            <p class="help is-danger">{errors.password}</p>
          )}
        </div>
        <br />
        <div class="field" className="is-flex mb-3">
          <input
            type="checkbox"
            id="rememberMe"
            checked={rememberMe}
            onChange={() => setRememberMe(!rememberMe)}
          />
          <label
            for="rememberMe"
            class="is-family-secondary is-size-7 has-text-grey-darker"
            style={{
              fontWeight: "700",
              marginLeft: "10px",
              //   textAlign: "center",
              //   display: "inline-block",
            }}
          >
            Remember Me
          </label>
        </div>
        <button
          class="button is-medium is-block is-fullwidth is-primary"
          type="button"
          onClick={onButtonClick}
        >
          Login <FontAwesomeIcon icon={faArrowRight} />
        </button>
      </form>
      <br />
      <nav class="level">
        <div class="level-item has-text-centered">
          <div>
            <Link to="/forgot-password" className="is-size-7-tablet">
              Forgot Password?
            </Link>
          </div>
        </div>
        <div class="level-item has-text-centered">
          <div>
            <Link to="/register-step-1" className="is-size-7-tablet is-hidden">
              Create an Account
            </Link>
          </div>
        </div>
      </nav>
    </div>
  );
};

export default LoginMobile;
