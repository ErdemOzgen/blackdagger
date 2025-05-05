import * as React from 'react';
import { styled, createTheme, ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import MuiDrawer from '@mui/material/Drawer';
import Box from '@mui/material/Box';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Typography from '@mui/material/Typography';
import IconButton from '@mui/material/IconButton';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import { mainListItems } from './menu';
import { Grid, MenuItem, Select, Divider } from '@mui/material';
import { AppBarContext } from './contexts/AppBarContext';
import { Link } from 'react-router-dom';
import blackdaggerImage from './assets/images/blackdagger.png';

const drawerWidth = 240;
const drawerWidthClosed = 64;

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme, open }) => ({
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(['width', 'margin'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  ...(open && {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  }),
}));

const Drawer = styled(MuiDrawer, {
  shouldForwardProp: (prop) => prop !== 'open',
})(({ theme, open }) => ({
  '& .MuiDrawer-paper': {
    backgroundColor: theme.palette.background.paper,
    position: 'relative',
    whiteSpace: 'nowrap',
    width: open ? drawerWidth : drawerWidthClosed,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.standard,
    }),
    overflowX: 'hidden',
    boxSizing: 'border-box',
  },
}));

const mdTheme = createTheme({
  palette: {
    mode: 'dark',
    background: {
      default: '#121212',
      paper: '#1e1e1e',
    },
    text: {
      primary: '#FFFEFE',
      secondary: '#CCCCCC',
    },
  },
  typography: {
    fontFamily: 'Inter',
  },
});

type DashboardContentProps = {
  title: string;
  navbarColor: string;
  version: string;
  children?: React.ReactNode;
};

function Content({ title, navbarColor, children }: DashboardContentProps) {
  const [scrolled, setScrolled] = React.useState(false);
  const [drawerOpen, setDrawerOpen] = React.useState(true);
  const containerRef = React.useRef<HTMLDivElement>(null);
  const gradientColor = navbarColor || '#171617';

  const toggleDrawer = () => {
    setDrawerOpen((prev) => !prev);
  };

  return (
    <ThemeProvider theme={mdTheme}>
      <Box sx={{ display: 'flex' }}>
        <CssBaseline />
        <AppBar position="absolute" open={drawerOpen}>
          <Toolbar
            sx={{
              pr: 2,
              display: 'flex',
              justifyContent: 'space-between',
            }}
          >
            <IconButton
              edge="start"
              color="inherit"
              onClick={toggleDrawer}
              sx={{ mr: 2 }}
            >
              {drawerOpen ? <ChevronLeftIcon /> : <MenuIcon />}
            </IconButton>

            <AppBarContext.Consumer>
              {(context) => (
                <NavBarTitleText visible={scrolled}>
                  {context.title}
                </NavBarTitleText>
              )}
            </AppBarContext.Consumer>

            <Box sx={{ display: 'flex', alignItems: 'center' }}>
              <Link
                to="/dashboard"
                style={{ textDecoration: 'none', marginRight: '10px' }}
              >
                <NavBarTitleText>{title || 'Blackdagger'}</NavBarTitleText>
              </Link>

              <AppBarContext.Consumer>
                {(context) =>
                  context.remoteNodes && context.remoteNodes.length > 0 ? (
                    <Select
                      sx={{
                        backgroundColor: 'white',
                        color: 'black',
                        borderRadius: '5px',
                        border: '1px solid #ccc',
                        ml: 2,
                        height: '30px',
                        width: '150px',
                      }}
                      value={context.selectedRemoteNode}
                      onChange={(e) => context.selectRemoteNode(e.target.value)}
                    >
                      {context.remoteNodes.map((node) => (
                        <MenuItem key={node} value={node}>
                          {node}
                        </MenuItem>
                      ))}
                    </Select>
                  ) : null
                }
              </AppBarContext.Consumer>
            </Box>
          </Toolbar>
        </AppBar>

        <Drawer variant="permanent" open={drawerOpen}>
          <Toolbar
            sx={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: drawerOpen ? 'center' : 'center',
              px: [1],
              py: 2,
            }}
          >
            <Box
              component="img"
              src={blackdaggerImage}
              alt="Blackdagger"
              sx={{
                width: drawerOpen ? 240 : 0,
                transition: 'width 0.4s ',
              }}
            />
          </Toolbar>

          <Divider />
          <Box
            sx={{
              background: `linear-gradient(0deg, ${mdTheme.palette.background.default} 0%, ${gradientColor} 100%)`,
              height: '100%',
            }}
          >
            <List component="nav" sx={{ pl: drawerOpen ? '12px' : '4px' }}>
              {mainListItems}
            </List>
          </Box>
        </Drawer>

        <Box
          component="main"
          sx={{
            flexGrow: 1,
            height: '100vh',
            overflow: 'auto',
            backgroundColor: '#171617',
          }}
        >
          <Toolbar />
          <Grid
            container
            ref={containerRef}
            onScroll={() => {
              const curr = containerRef.current;
              if (curr) {
                setScrolled(curr.scrollTop > 54);
              }
            }}
            sx={{ flex: 1, pb: 4, px: 3 }}
          >
            {children}
          </Grid>
        </Box>
      </Box>
    </ThemeProvider>
  );
}

type NavBarTitleTextProps = {
  children: string;
  visible?: boolean;
};

const NavBarTitleText = ({
  children,
  visible = true,
}: NavBarTitleTextProps) => (
  <Typography
    component="h1"
    variant="h6"
    sx={{
      fontWeight: '800',
      color: '#FFFEFE',
      opacity: visible ? 1 : 0,
      transition: 'opacity 0.2s',
    }}
  >
    {children}
  </Typography>
);

type DashboardProps = DashboardContentProps;

export default function Layout({ children, ...props }: DashboardProps) {
  return <Content {...props}>{children}</Content>;
}
