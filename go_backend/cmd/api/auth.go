package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	Issuer        string        //who is issueing this token
	Audience      string        //who should be able to use these token
	Secret        string        //strong secret key which we use to sign our tokens
	TokenExpiry   time.Duration //how lonh should our token last
	RefreshExpiry time.Duration // this is technically not a partof jwt standards but most people that use that tend to use refresh token as well in other words that jwt it self has short expiry time 15min/5min  but we also give the user refresh token which also use to authenticate user again back to the website that may last upto 2weeks in some case
	CookieDomain  string        //we will give refresh token to the user as a cookie , as a http only secure cookie which will not is accesible from java scrip but which will be included when me make request to the backend
	CookiePath    string        // which we will set to the "/"  the root level of our application
	CookieName    string        // name of the cookie

}

// This Type will only contain that is needed to issue a token

//will have a user in a database it will be associated with those user and store or collect their information and assign them a jwt token
type jwtUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type TokenPairs struct {
	Token        string `json:"access_token"`  // will use this token to provide access to the user while communicating to the backend
	RefreshToken string `json:"refresh_token"` //We provide refresh tokens as a way to obtain a new access token when the current access token has expired. Access tokens are often used in authentication systems to grant access to protected resources, such as APIs or web applications.
	// Access tokens have a limited lifespan, typically ranging from a few minutes to a few hours, to minimize the risk of unauthorized access if a token is stolen or compromised. When an access token expires, the user would normally have to log in again to obtain a new one. However, this can be inconvenient and disrupt the user experience.
	// A refresh token allows the user to obtain a new access token without having to log in again. The refresh token is typically a long-lived token that can be used to obtain a new access token when the current one has expired. The refresh token is typically more secure than the access token, as it is usually only used for obtaining new access tokens, and is not directly used to access protected resources.
	// Providing a refresh token is a common way to balance security and user experience in authentication systems, as it allows users to access protected resources without the inconvenience of logging in again each time their access token expires.
}

// every time you have a jwt issued that jwt has certain things called claim 
//you might claim that this jwt token is only for certain audience
// there are things you have to have in there and there are things you cannot have in there this logic explains claims  that only allowed user are able to access it
type Claims struct{
	jwt.RegisteredClaims
}

func (j *Auth) GenerateTokenPair(user *jwtUser)(TokenPairs,error){
	//create a token
	token:= jwt.New(jwt.SigningMethodHS256)

	// set claims
	Claims:=token.Claims.(jwt.MapClaims)
	Claims["name"]=fmt.Sprintf("%s,%s",user.FirstName,user.LastName)
	Claims["sub"]=fmt.Sprint(user.ID)
	Claims["aud"]=j.Audience
	Claims["iss"]=j.Issuer
	Claims["iat"]=time.Now().UTC().Unix()
	Claims["jwt"]="JWT"

	//set the expiry for JWT which be shorter than refresh token
	Claims["exp"]=time.Now().UTC().Add(j.TokenExpiry).Unix()

	//Create a signed Token
	signedAccessToken,err:=token.SignedString([]byte(j.Secret))
	if err !=nil{
		return TokenPairs{},err
	}

	//Create a ResfreshToken and set Claims
	refreshToken:=jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims:=refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"]=(user.ID)
	refreshTokenClaims["iat"]=time.Now().UTC().Unix()


	//Set the Expiry for RefreshTokens
	refreshTokenClaims["exp"]=time.Now().UTC().Add(j.RefreshExpiry).Unix()

	//Create a signed RefreshToken
	signedrefreshtoken,err:=refreshToken.SignedString([]byte(j.Secret))
	if err !=nil{
		return TokenPairs{},err
	}	

	//Create Token Pairs(Variable of type token pair) and Populate signed tokens
		var tokenPairs = TokenPairs{
			Token: signedAccessToken,
			RefreshToken: signedrefreshtoken,
		}
	//Return TokenPairs
	return tokenPairs, nil

}

/* in general, it's common to create a function like getRefreshCookie to generate a new HTTP cookie containing a JWT refresh token.
A JWT refresh token is a type of token that is used to authenticate a user and obtain a new access token after the previous one has expired. Typically, the refresh token is sent to the server in an HTTP cookie, which is then used to generate a new access token. The getRefreshCookie function would be responsible for creating and returning this cookie.
By passing the refreshToken as an input parameter to the function, the function can generate a new cookie with the refreshToken value set as the cookie value. This allows the server to read the refreshToken value from the cookie and use it to generate a new access token when necessary.
Overall, creating a function like getRefreshCookie is a common pattern when working with JWT refresh tokens, as it provides a standardized way to generate the necessary HTTP cookies.



*/

func (j *Auth) GetRefreshCookie(refreshToken string)*http.Cookie{
	return &http.Cookie{
		Name: j.CookieName,
		Path: j.CookiePath,
		Value: refreshToken,
		Expires: time.Now().Add(j.RefreshExpiry),
		MaxAge: int(j.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain: j.CookieDomain,
		HttpOnly: true,
		Secure: true,
	}
}

func (j *Auth) GetExpiredRefreshCookie(refreshToken string)*http.Cookie{
	return &http.Cookie{
		Name: j.CookieName,
		Path: j.CookiePath,
		Value: "",
		Expires: time.Unix(0,0),
		MaxAge: -1,
		SameSite: http.SameSiteStrictMode,
		Domain: j.CookieDomain,
		HttpOnly: true,
		Secure: true,
	}
}