import React, {Component} from "react";

import axios from "axios"
import {Card, Header, Form, Input,Icon} from "semantic-ui-react"

let endPoint = "http://localhost:9000"

class ToDoList extends Component {
    constructor(props) {
        super(props);

        this.state = {
            task:"",
            items:[],
        }
    }

    componentDidMount() {
        this.getTask();
    }

    onChange = (event)=> {
        this.setState({
            [event.target.name]: event.target.value
        })
    }

    onSubmit = () => {
        let {task} = this.state;
        if (task) {
            axios.post(endPoint+"/api/task",
                {task,},
                {headers: {"Content-Type": "application/x-www-form-urlencoded",}
                }).then((res)=>{
                    this.getTask();
                    this.setState({
                        task:"",
                    });
                    console.log(res);
            })
        }
    }

    getTask = () => {
        axios.get(endPoint + "/api/task").then((res)=>{
            if (res.data) {
                this.setState({
                    items: res.data.map((item)=>{
                        let color = "yellow";
                        let style = {
                            wordWrap: "break-word",
                        };

                        if(item.status) {
                            color="green"
                            style["texeDecorationLine"] = "line-through"
                        }
                        return(
                            <Card key={item._id} color={color} fluid className="rough">
                                <Card.Content textAlign="left">
                                    <div style={style}>{item.task}</div>
                                </Card.Content>
                                <Card.Meta textAlign="right">
                                    <Icon
                                        name="check circle"
                                        color="blue"
                                        onClick={()=>this.updateTask(item._id)}>
                                    </Icon>
                                    <span style={{paddingRight:10}}>Undo</span>

                                    <Icon
                                        name="delete"
                                        color="red"
                                        onClick={()=>this.deleteTask(item._id)}>
                                    </Icon>
                                    <span style={{paddingRight:10}}>Delete</span>

                                </Card.Meta>
                            </Card>
                        );
                    }),
                });
            }else {
                this.setState({
                    items: [],
                })
            }
        })
    }

    undoTask = (id) => {
        axios.put(endPoint+"/api/undoTask" + id, {
            headers:{
                "Content-Type": "application/x-www-form-urlencoded",
            },
        }).then((res)=>{
            console.log(res);
            this.getTask()
        })
    }

    updateTask = (id) => {
        axios.put(endPoint+"/api/task" + id, {
            headers:{
                "Content-Type": "application/x-www-form-urlencoded",
            },
        }).then((res)=>{
            console.log(res);
            this.getTask()
        })
    }

    deleteTask = (id) => {
        axios.put(endPoint+"/api/deleteTask" + id, {
            headers:{
                "Content-Type": "application/x-www-form-urlencoded",
            },
        }).then((res)=>{
            console.log(res);
            this.getTask()
        })
    }

    render() {
        return(
          <div>
              <div className="row">
                  <Header className="header" as="h2" color="yellow">
                    TO DO LIST
                  </Header>
              </div>
              <div className="row">
                  <Form onSubmit={this.onSubmit}>
                      <Input type="text" name="task" onChange={this.onChange} value={this.state.task} placeholder="Create Task">
                      </Input>

                      {/*<Button> Create Task </Button>*/}
                  </Form>
              </div>

              <div className="row">
                  <Card.Group>{this.state.items}</Card.Group>
              </div>
          </div>
        );
    }
}

export default ToDoList;






