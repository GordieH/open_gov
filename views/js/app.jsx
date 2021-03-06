const AUTH0_CLIENT_ID = "9NH1MWLWcM54FRBn0Xvo2dOEFaJKG3gr";
const AUTH0_DOMAIN = "dev-vb3a8shg.us.auth0.com";
const AUTH0_CALLBACK_URL = location.href;
const AUTH0_API_AUDIENCE = "https://gordieh.opengov/";

class App extends React.Component {
  parseHash() {
    this.auth0 = new auth0.WebAuth({
      domain: AUTH0_DOMAIN,
      clientID: AUTH0_CLIENT_ID,
    });
    this.auth0.parseHash(window.location.hash, (err, authResult) => {
      if (err) {
        return console.log(err);
      }
      if (
        authResult !== null &&
        authResult.accessToken !== null &&
        authResult.idToken !== null
      ) {
        localStorage.setItem("access_token", authResult.accessToken);
        localStorage.setItem("id_token", authResult.idToken);
        localStorage.setItem(
          "profile",
          JSON.stringify(authResult.idTokenPayload)
        );
        window.location = window.location.href.substr(
          0,
          window.location.href.indexOf("#")
        );
      }
    });
  }

  setup() {
    $.ajaxSetup({
      beforeSend: (r) => {
        if (localStorage.getItem("access_token")) {
          r.setRequestHeader(
            "Authorization",
            "Bearer " + localStorage.getItem("access_token")
          );
        }
      },
    });
  }

  setState() {
    let idToken = localStorage.getItem("id_token");
    if (idToken) {
      this.loggedIn = true;
    } else {
      this.loggedIn = false;
    }
  }

  componentWillMount() {
    this.setup();
    this.parseHash();
    this.setState();
  }

  render() {
    if (this.loggedIn) {
      return <LoginHome />;
    }
    return <Home />;
  }
}

class Home extends React.Component {
  constructor(props) {
    super(props);
    this.authenticate = this.authenticate.bind(this);
  }

  authenticate() {
    this.WebAuth = new auth0.WebAuth({
      domain: AUTH0_DOMAIN,
      clientID: AUTH0_CLIENT_ID,
      scope: "openid profile",
      audience: AUTH0_API_AUDIENCE,
      responseType: "token id_token",
      redirectUri: AUTH0_CALLBACK_URL,
    });
    this.WebAuth.authorize();
  }

  render() {
    return (
      <div className="container">
        <div className="row">
          <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
            <h1>Open-Gov</h1>
            <p>An open-source app for engaging citizens with government</p>
            <p>Sign in to get started </p>
            <a
              onClick={this.authenticate}
              className="btn btn-primary btn-lg btn-login btn-block"
            >
              Sign In
            </a>
          </div>
        </div>
      </div>
    );
  }
}

class LoginHome extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      user_guid: null,
      reps: [],
    };
  }

  logout = () => {
    localStorage.removeItem("id_token");
    localStorage.removeItem("access_token");
    localStorage.removeItem("profile");
    location.reload();
  };

  serverRequest = () => {
    $.get("http://localhost:3000/api/localreps", (res) => {
      this.setState({
        user_guid: res.user_guid,
        reps: res.users_rep_list,
      });
    });
  };

  handleAdd = (repName) => {
    console.log('repName: ', repName)
    // TODO: Add API call here to get all of the rep's information
  };

  handleDelete = (guid) => {
    // TODO: I don't think we need to make a network call here, depends on
    // if we want to store a user's modified reps in a database or in local state
    $.post(`http://localhost:3000/api/localreps/edit?editTask=remove&user_guid=55ee03f2dcd8c8e46b91cbb2e70d9e&rep_guid=${guid}`, res => {
       return res;
    });
    const updatedReps = this.state.reps.filter((rep) => rep.guid !== guid);
    this.setState({
      reps: updatedReps,
    });
  };

  componentDidMount() {
    this.serverRequest();
  }

  render() {
    const userList = this.state.reps;
    return (
      <div className="container">
        <br />
        <span className="pull-right">
          <a onClick={this.logout}>Log out</a>
        </span>
        <h2>Open-Gov</h2>
        <p>Hey user</p>
        <RepNameForm addRep={this.handleAdd} />
        <br />
        <div className="row">
          <div className="container">
            {userList &&
              userList.map((localRep, i) => {
                return (
                  <RepCard
                    key={i}
                    localRep={localRep}
                    deleteRep={this.handleDelete}
                  />
                );
              })}
          </div>
        </div>
      </div>
    );
  }
}
class RepNameForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      input: '',
    }
  }

  handleChange = (event) => {
    this.setState({
      input: event.target.value,
    });
  }

  handleSubmit = () => {
    this.props.addRep(this.state.input);
  }

  render() {
    return (
      <form onSubmit={() => this.handleSubmit()}>
        <label>
          Enter a representative's name to add: 
          <input type="text" value={this.state.value} onChange={this.handleChange} />
        </label>
        <input type="submit" value="Submit" />
      </form>
    );
  }
}

class RepCard extends React.Component {
  constructor(props) {
    super(props);
  }

  deleteRep = (id) => {
    this.props.deleteRep(id);
  };

  render() {
    const { localRep } = this.props;

    return (
      <div className="col-xs-4">
        <div className="panel panel-default">
          <div className="panel-heading">
            {localRep.name} {localRep.LastName}{" "}
            <span className="pull-right"></span>
          </div>
          <div className="panel-body joke-hld">Office: {localRep.office}</div>
          <div className="panel-body joke-hld">
            Location: {localRep.location}
          </div>
          <div className="panel-body joke-hld">
            Percent Votes Missed: {localRep.percent_missed_votes}%
          </div>
          <div className="panel-body joke-hld">
            Percent Votes With Party: {localRep.percent_votes_with_party}%
          </div>
          <div className="panel-body joke-hld">
            <a href={localRep.gov_web}>Goverment Web Page</a>
            <div> </div>
            <a href={`https://www.twitter.com/${localRep.twitter}`}>Twitter</a>
          </div>
        </div>
        <button type="button" onClick={() => this.deleteRep(localRep.guid)}>
          Remove Rep
        </button>
      </div>
    );
  }
}

ReactDOM.render(<App />, document.getElementById("app"));
